package repository

import (
	"sync"
	"time"

	"cyskillswap/internal/constants"
	bizerrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// NeedQueueStore 求助与候选名单的内存存储，所有状态迁移在单把锁下完成，保证
// “同时只有一个约谈”“退出后自动接手”等约束在并发下也成立。
type NeedQueueStore struct {
	mu             sync.RWMutex
	needs          []model.Need
	responses      []model.NeedResponse
	nextResponseID int
	nextSeq        int
}

var needQueue = newNeedQueueStore()

func newNeedQueueStore() *NeedQueueStore {
	s := &NeedQueueStore{nextResponseID: 1, nextSeq: 1}
	s.seed()
	return s
}

// ---------- 查询（返回拷贝，避免外部修改内部状态） ----------

// ListNeeds 返回全部求助，Responses 为候选名单中的活跃人数（约谈中+排队中）
func ListNeeds() []model.Need {
	needQueue.mu.RLock()
	defer needQueue.mu.RUnlock()
	needs := make([]model.Need, len(needQueue.needs))
	for i, need := range needQueue.needs {
		need.Responses = countActiveLocked(need.ID)
		needs[i] = need
	}
	return needs
}

// NeedSnapshot 返回单条求助及其候选名单（按提交先后排序）的拷贝
func NeedSnapshot(needID int) (model.Need, []model.NeedResponse, bool) {
	needQueue.mu.RLock()
	defer needQueue.mu.RUnlock()
	need, ok := findNeedLocked(needID)
	if !ok {
		return model.Need{}, nil, false
	}
	need.Responses = countActiveLocked(needID)
	return need, copyResponsesLocked(needID), true
}

// ---------- 状态迁移 ----------

// SubmitResponse 响应者提交候选：能帮的方式 + 可上门时段，按提交先后入队
func SubmitResponse(needID int, input model.SubmitResponseInput) (model.NeedResponse, error) {
	needQueue.mu.Lock()
	defer needQueue.mu.Unlock()
	need, ok := findNeedLocked(needID)
	if !ok {
		return model.NeedResponse{}, bizerrors.ErrNeedNotFound
	}
	if need.Status == constants.NeedStatusClosed {
		return model.NeedResponse{}, bizerrors.ErrNeedClosed
	}
	if input.Responder == need.Requester {
		return model.NeedResponse{}, bizerrors.ErrRespondOwnNeed
	}
	for _, resp := range needQueue.responses {
		if resp.NeedID == needID && resp.Responder == input.Responder && isActiveStatus(resp.Status) {
			return model.NeedResponse{}, bizerrors.ErrDuplicateResponse
		}
	}
	resp := model.NeedResponse{
		ID:          needQueue.nextResponseID,
		NeedID:      needID,
		Responder:   input.Responder,
		HelpOffer:   input.HelpOffer,
		VisitSlots:  append([]string(nil), input.VisitSlots...),
		Status:      constants.ResponseStatusWaiting,
		Seq:         needQueue.nextSeq,
		SubmittedAt: time.Now().Format(constants.DateTimeLayout),
	}
	needQueue.nextResponseID++
	needQueue.nextSeq++
	needQueue.responses = append(needQueue.responses, resp)
	return resp, nil
}

// PickInterview 发布人从候选名单中挑一位约谈；同一求助同时只能有一个约谈
func PickInterview(needID, responseID int, operator string) error {
	needQueue.mu.Lock()
	defer needQueue.mu.Unlock()
	need, resp, err := lockNeedAndResponse(needID, responseID)
	if err != nil {
		return err
	}
	if operator != need.Requester {
		return bizerrors.ErrNotPublisher
	}
	if need.Status == constants.NeedStatusClosed {
		return bizerrors.ErrNeedClosed
	}
	if resp.Status != constants.ResponseStatusWaiting {
		return bizerrors.ErrCandidateNotWaiting
	}
	if _, interviewing := findInterviewingLocked(needID); interviewing {
		return bizerrors.ErrInterviewOccupied
	}
	resp.Status = constants.ResponseStatusInterviewing
	return nil
}

// WithdrawResponse 响应者退出轮候；若退出的是当前约谈人，队首自动接手约谈
func WithdrawResponse(needID, responseID int, operator string) error {
	needQueue.mu.Lock()
	defer needQueue.mu.Unlock()
	need, resp, err := lockNeedAndResponse(needID, responseID)
	if err != nil {
		return err
	}
	if operator != resp.Responder {
		return bizerrors.ErrNotCandidate
	}
	if need.Status == constants.NeedStatusClosed {
		return bizerrors.ErrNeedClosed
	}
	if resp.Status == constants.ResponseStatusWithdrawn {
		return bizerrors.ErrAlreadyWithdrawn
	}
	wasInterviewing := resp.Status == constants.ResponseStatusInterviewing
	resp.Status = constants.ResponseStatusWithdrawn
	if wasInterviewing {
		promoteNextLocked(needID)
	}
	return nil
}

// CloseNeed 发布人关闭求助，候选名单停止（不再接受响应/约谈/退出）
func CloseNeed(needID int, operator string) error {
	needQueue.mu.Lock()
	defer needQueue.mu.Unlock()
	need, ok := findNeedLocked(needID)
	if !ok {
		return bizerrors.ErrNeedNotFound
	}
	if operator != need.Requester {
		return bizerrors.ErrNotPublisher
	}
	if need.Status == constants.NeedStatusClosed {
		return bizerrors.ErrNeedAlreadyClosed
	}
	for i := range needQueue.needs {
		if needQueue.needs[i].ID == needID {
			needQueue.needs[i].Status = constants.NeedStatusClosed
			break
		}
	}
	return nil
}

// ---------- 内部工具（调用方须已持有锁） ----------

func findNeedPtrLocked(needID int) *model.Need {
	for i := range needQueue.needs {
		if needQueue.needs[i].ID == needID {
			return &needQueue.needs[i]
		}
	}
	return nil
}

func findNeedLocked(needID int) (model.Need, bool) {
	need := findNeedPtrLocked(needID)
	if need == nil {
		return model.Need{}, false
	}
	return *need, true
}

func lockNeedAndResponse(needID, responseID int) (model.Need, *model.NeedResponse, error) {
	need, ok := findNeedLocked(needID)
	if !ok {
		return model.Need{}, nil, bizerrors.ErrNeedNotFound
	}
	for i := range needQueue.responses {
		resp := &needQueue.responses[i]
		if resp.ID == responseID && resp.NeedID == needID {
			return need, resp, nil
		}
	}
	return model.Need{}, nil, bizerrors.ErrResponseNotFound
}

func isActiveStatus(status string) bool {
	return status == constants.ResponseStatusWaiting || status == constants.ResponseStatusInterviewing
}

func countActiveLocked(needID int) int {
	count := 0
	for _, resp := range needQueue.responses {
		if resp.NeedID == needID && isActiveStatus(resp.Status) {
			count++
		}
	}
	return count
}

func copyResponsesLocked(needID int) []model.NeedResponse {
	var result []model.NeedResponse
	for _, resp := range needQueue.responses {
		if resp.NeedID != needID {
			continue
		}
		resp.VisitSlots = append([]string(nil), resp.VisitSlots...)
		result = append(result, resp)
	}
	return result
}

func findInterviewingLocked(needID int) (*model.NeedResponse, bool) {
	for i := range needQueue.responses {
		resp := &needQueue.responses[i]
		if resp.NeedID == needID && resp.Status == constants.ResponseStatusInterviewing {
			return resp, true
		}
	}
	return nil, false
}

// promoteNextLocked 约谈对象退出后，按提交先后让队首的排队者自动接手
func promoteNextLocked(needID int) {
	var next *model.NeedResponse
	for i := range needQueue.responses {
		resp := &needQueue.responses[i]
		if resp.NeedID == needID && resp.Status == constants.ResponseStatusWaiting {
			if next == nil || resp.Seq < next.Seq {
				next = resp
			}
		}
	}
	if next != nil {
		next.Status = constants.ResponseStatusInterviewing
	}
}
