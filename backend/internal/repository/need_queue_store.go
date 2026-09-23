package repository

import (
	"sort"
	"sync"
	"time"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// needQueueStore 内存版求助队列存储，负责名单状态的原子流转。
type needQueueStore struct {
	mu         sync.RWMutex
	needs      map[int]*model.Need
	responses  map[int][]*model.NeedResponse // needID → 按提交先后排序的候选名单
	nextRespID int
}

var queueStore = seedNeedQueueStore()

func seedNeedQueueStore() *needQueueStore {
	s := &needQueueStore{
		needs:      map[int]*model.Need{},
		responses:  map[int][]*model.NeedResponse{},
		nextRespID: 1,
	}
	needs := []model.Need{
		{ID: 1, Requester: "孟野", Title: "找人帮忙拍乐队宣传照", Category: "摄影", Campus: "西校区", ExpectTime: "本周六上午", BudgetType: "技能交换", Description: "可交换 3 次吉他课，希望会调色和室外构图。", Status: constants.NeedStatusOpen},
		{ID: 2, Requester: "许安", Title: "求教 Python 数据分析", Category: "编程", Campus: "中心校区", ExpectTime: "周二晚", BudgetType: "小额报酬", Description: "论文问卷数据需要清洗和画图，最好有 pandas 经验。", Status: constants.NeedStatusOpen},
		{ID: 3, Requester: "林澈", Title: "想学吉他扫弦入门", Category: "乐器", Campus: "东校区", ExpectTime: "周三晚", BudgetType: "技能交换", Description: "用摄影课交换吉他基础，希望同校区或线上。", Status: constants.NeedStatusOpen},
		{ID: 4, Requester: "周芮", Title: "找人模拟英语求职面试", Category: "外语", Campus: "中心校区", ExpectTime: "周四晚", BudgetType: "请吃饭", Description: "下周外企实习面试，想练自我介绍和问答环节。", Status: constants.NeedStatusClosed},
	}
	for i := range needs {
		n := needs[i]
		s.needs[n.ID] = &n
	}
	seed := func(needID int, responder, helpMethod string, slots []string, status, submittedAt string) {
		s.appendResponseLocked(needID, responder, helpMethod, slots, status, submittedAt)
	}
	// 求助 1：林澈（演示账号）正在约谈中，其余两人排队
	seed(1, "林澈", "室外人像跟拍 + 精修 9 张，可带补光灯", []string{"周六上午", "周日下午"}, constants.ResponseStatusInterviewing, "2026-09-19 09:12")
	seed(1, "周芮", "现场花絮记录，可带反光板协助打光", []string{"周六上午"}, constants.ResponseStatusWaiting, "2026-09-19 11:40")
	seed(1, "许安", "活动跟拍经验丰富，当天可出快修图", []string{"周六上午", "周六下午"}, constants.ResponseStatusWaiting, "2026-09-20 08:55")
	// 求助 2：第一位约谈后撤回，周芮自动接手，林澈排队等待
	seed(2, "孟野", "可教 Excel 透视表过渡方案", []string{"周二晚"}, constants.ResponseStatusWithdrawn, "2026-09-18 14:03")
	seed(2, "周芮", "pandas 清洗 + 可视化模板演示 1 小时", []string{"周二晚"}, constants.ResponseStatusInterviewing, "2026-09-18 16:27")
	seed(2, "林澈", "可帮忙整理问卷数据并讲解绘图思路", []string{"周二晚", "周四晚"}, constants.ResponseStatusWaiting, "2026-09-19 20:11")
	// 求助 3：演示账号是发布人，周芮已退出，两人排队待挑
	seed(3, "孟野", "民谣扫弦入门 2 次课，想换摄影课", []string{"周三晚", "周六上午"}, constants.ResponseStatusWaiting, "2026-09-19 10:02")
	seed(3, "周芮", "可教基础和弦转换", []string{"周五晚"}, constants.ResponseStatusExited, "2026-09-19 15:46")
	seed(3, "许安", "尤克里里转吉他，可带练扫弦节奏", []string{"周五晚", "周日下午"}, constants.ResponseStatusWaiting, "2026-09-20 09:31")
	// 求助 4：已关闭，名单停止
	seed(4, "林澈", "可模拟英文面试问答并纠正发音", []string{"周四晚"}, constants.ResponseStatusClosed, "2026-09-15 19:20")
	seed(4, "孟野", "有外企实习面试经验，可陪练自我介绍", []string{"周四晚", "周五晚"}, constants.ResponseStatusClosed, "2026-09-16 08:44")
	s.refreshResponseCountsLocked()
	return s
}

// appendResponseLocked 追加候选条目，位次按提交先后递增。调用前须持有写锁或在种子阶段。
func (s *needQueueStore) appendResponseLocked(needID int, responder, helpMethod string, slots []string, status, submittedAt string) *model.NeedResponse {
	entry := &model.NeedResponse{
		ID:          s.nextRespID,
		NeedID:      needID,
		Responder:   responder,
		HelpMethod:  helpMethod,
		TimeSlots:   slots,
		Status:      status,
		Position:    len(s.responses[needID]) + 1,
		SubmittedAt: submittedAt,
	}
	s.nextRespID++
	s.responses[needID] = append(s.responses[needID], entry)
	return entry
}

func (s *needQueueStore) refreshResponseCountsLocked() {
	for id, need := range s.needs {
		need.Responses = len(s.responses[id])
	}
}

func cloneNeed(n *model.Need) model.Need { return *n }

func cloneResponse(r *model.NeedResponse) model.NeedResponse {
	c := *r
	c.TimeSlots = append([]string(nil), r.TimeSlots...)
	return c
}

// ListNeeds 按 ID 顺序返回全部求助。
func (s *needQueueStore) ListNeeds() []model.Need {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]int, 0, len(s.needs))
	for id := range s.needs {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	needs := make([]model.Need, 0, len(ids))
	for _, id := range ids {
		needs = append(needs, cloneNeed(s.needs[id]))
	}
	return needs
}

// GetNeed 返回单个求助。
func (s *needQueueStore) GetNeed(id int) (model.Need, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	need, ok := s.needs[id]
	if !ok {
		return model.Need{}, false
	}
	return cloneNeed(need), true
}

// ListResponses 返回某求助按提交先后排序的候选名单。
func (s *needQueueStore) ListResponses(needID int) []model.NeedResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := s.responses[needID]
	out := make([]model.NeedResponse, 0, len(entries))
	for _, entry := range entries {
		out = append(out, cloneResponse(entry))
	}
	return out
}

// SubmitResponse 响应者提交能帮的方式与可上门时段，加入队尾。
func (s *needQueueStore) SubmitResponse(needID int, operator, helpMethod string, slots []string) (model.NeedResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	need, ok := s.needs[needID]
	if !ok {
		return model.NeedResponse{}, errors.New(constants.ErrCodeNeedNotFound, constants.ErrMsgNeedNotFound)
	}
	if need.Status != constants.NeedStatusOpen {
		return model.NeedResponse{}, errors.New(constants.ErrCodeNeedClosed, constants.ErrMsgNeedClosed)
	}
	if need.Requester == operator {
		return model.NeedResponse{}, errors.New(constants.ErrCodeSelfResponse, constants.ErrMsgSelfResponse)
	}
	for _, entry := range s.responses[needID] {
		if entry.Responder == operator && isActiveResponse(entry.Status) {
			return model.NeedResponse{}, errors.New(constants.ErrCodeDuplicateResponse, constants.ErrMsgDuplicateResponse)
		}
	}
	entry := s.appendResponseLocked(needID, operator, helpMethod, slots, constants.ResponseStatusWaiting, time.Now().Format("2006-01-02 15:04"))
	s.refreshResponseCountsLocked()
	return cloneResponse(entry), nil
}

// PickCandidate 发布人从排队名单中挑一位约谈；同一求助同时只能有一位约谈中。
func (s *needQueueStore) PickCandidate(needID, responseID int, operator string) (model.NeedResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	need, ok := s.needs[needID]
	if !ok {
		return model.NeedResponse{}, errors.New(constants.ErrCodeNeedNotFound, constants.ErrMsgNeedNotFound)
	}
	if need.Requester != operator {
		return model.NeedResponse{}, errors.New(constants.ErrCodeNotRequester, constants.ErrMsgNotRequester)
	}
	if need.Status != constants.NeedStatusOpen {
		return model.NeedResponse{}, errors.New(constants.ErrCodeNeedClosed, constants.ErrMsgNeedClosed)
	}
	target := findResponseLocked(s.responses[needID], responseID)
	if target == nil {
		return model.NeedResponse{}, errors.New(constants.ErrCodeResponseNotFound, constants.ErrMsgResponseNotFound)
	}
	if target.Status != constants.ResponseStatusWaiting {
		return model.NeedResponse{}, errors.New(constants.ErrCodeCandidateNotWaiting, constants.ErrMsgCandidateNotWaiting)
	}
	for _, entry := range s.responses[needID] {
		if entry.Status == constants.ResponseStatusInterviewing {
			return model.NeedResponse{}, errors.New(constants.ErrCodeInterviewInProgress, constants.ErrMsgInterviewInProgress)
		}
	}
	target.Status = constants.ResponseStatusInterviewing
	return cloneResponse(target), nil
}

// ExitResponse 响应者退出：排队中记为已退出，约谈中记为已撤回并由下一位自动接手。
// 返回被更新的条目与自动接手的条目（无则 nil）。
func (s *needQueueStore) ExitResponse(needID int, operator string) (model.NeedResponse, *model.NeedResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.needs[needID]; !ok {
		return model.NeedResponse{}, nil, errors.New(constants.ErrCodeNeedNotFound, constants.ErrMsgNeedNotFound)
	}
	var mine *model.NeedResponse
	for _, entry := range s.responses[needID] {
		if entry.Responder == operator && isActiveResponse(entry.Status) {
			mine = entry
			break
		}
	}
	if mine == nil {
		return model.NeedResponse{}, nil, errors.New(constants.ErrCodeResponseNotActive, constants.ErrMsgResponseNotActive)
	}
	var promoted *model.NeedResponse
	if mine.Status == constants.ResponseStatusInterviewing {
		mine.Status = constants.ResponseStatusWithdrawn
		promoted = promoteNextLocked(s.responses[needID])
	} else {
		mine.Status = constants.ResponseStatusExited
	}
	return cloneResponse(mine), promoted, nil
}

// CloseNeed 发布人关闭求助：名单停止变动，排队中的候选人记为轮候结束。
func (s *needQueueStore) CloseNeed(needID int, operator string) (model.Need, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	need, ok := s.needs[needID]
	if !ok {
		return model.Need{}, errors.New(constants.ErrCodeNeedNotFound, constants.ErrMsgNeedNotFound)
	}
	if need.Requester != operator {
		return model.Need{}, errors.New(constants.ErrCodeNotRequester, constants.ErrMsgNotRequester)
	}
	if need.Status != constants.NeedStatusOpen {
		return model.Need{}, errors.New(constants.ErrCodeNeedClosed, constants.ErrMsgNeedClosed)
	}
	need.Status = constants.NeedStatusClosed
	for _, entry := range s.responses[needID] {
		if entry.Status == constants.ResponseStatusWaiting {
			entry.Status = constants.ResponseStatusClosed
		}
	}
	return cloneNeed(need), nil
}

// promoteNextLocked 约谈对象撤回后，队首的排队者自动接手。
func promoteNextLocked(entries []*model.NeedResponse) *model.NeedResponse {
	for _, entry := range entries {
		if entry.Status == constants.ResponseStatusWaiting {
			entry.Status = constants.ResponseStatusInterviewing
			c := cloneResponse(entry)
			return &c
		}
	}
	return nil
}

func findResponseLocked(entries []*model.NeedResponse, id int) *model.NeedResponse {
	for _, entry := range entries {
		if entry.ID == id {
			return entry
		}
	}
	return nil
}

// isActiveResponse 仍在队列中的状态（排队中 / 约谈中）。
func isActiveResponse(status string) bool {
	return status == constants.ResponseStatusWaiting || status == constants.ResponseStatusInterviewing
}

// ---- 包级门面：供 service 层调用的求助队列操作 ----

// ListResponses 返回某求助按提交先后排序的候选名单。
func ListResponses(needID int) []model.NeedResponse { return queueStore.ListResponses(needID) }

// GetNeed 返回单个求助。
func GetNeed(id int) (model.Need, bool) { return queueStore.GetNeed(id) }

// SubmitResponse 响应者提交响应并加入队尾。
func SubmitResponse(needID int, operator, helpMethod string, slots []string) (model.NeedResponse, error) {
	return queueStore.SubmitResponse(needID, operator, helpMethod, slots)
}

// PickCandidate 发布人挑一位排队者约谈。
func PickCandidate(needID, responseID int, operator string) (model.NeedResponse, error) {
	return queueStore.PickCandidate(needID, responseID, operator)
}

// ExitResponse 响应者退出；约谈对象退出时由下一位自动接手。
func ExitResponse(needID int, operator string) (model.NeedResponse, *model.NeedResponse, error) {
	return queueStore.ExitResponse(needID, operator)
}

// CloseNeed 发布人关闭求助，名单停止。
func CloseNeed(needID int, operator string) (model.Need, error) {
	return queueStore.CloseNeed(needID, operator)
}
