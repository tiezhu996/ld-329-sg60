package service

import (
	"fmt"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/logger"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// resolveViewer 未指定用户时使用演示账号视角。
func resolveViewer(viewer string) string {
	if viewer == "" {
		return constants.DefaultViewer
	}
	return viewer
}

// ListNeedSummaries 求助列表：每条含当前约谈人与查看者自己的排队信息。
func ListNeedSummaries(viewer string) []model.NeedSummary {
	viewer = resolveViewer(viewer)
	needs := repository.ListNeeds()
	summaries := make([]model.NeedSummary, 0, len(needs))
	for _, need := range needs {
		summaries = append(summaries, buildNeedSummary(need, repository.ListResponses(need.ID), viewer))
	}
	return summaries
}

// GetNeedDetail 求助详情：完整候选名单与查看者可执行的操作。
func GetNeedDetail(needID int, viewer string) (model.NeedDetail, error) {
	viewer = resolveViewer(viewer)
	need, ok := repository.GetNeed(needID)
	if !ok {
		return model.NeedDetail{}, errors.New(constants.ErrCodeNeedNotFound, constants.ErrMsgNeedNotFound)
	}
	queue := repository.ListResponses(needID)
	summary := buildNeedSummary(need, queue, viewer)
	detail := model.NeedDetail{NeedSummary: summary, Queue: queue}
	isRequester := need.Requester == viewer
	needOpen := need.Status == constants.NeedStatusOpen
	detail.CanClose = isRequester && needOpen
	detail.CanPick = isRequester && needOpen && summary.CurrentInterview == ""
	detail.CanRespond = !isRequester && needOpen && summary.MyQueue == nil
	return detail, nil
}

// SubmitNeedResponse 响应者提交能帮的方式与可上门时段，按提交先后入队。
func SubmitNeedResponse(needID int, input model.SubmitResponseInput) (model.NeedDetail, error) {
	operator := resolveViewer(input.Operator)
	if input.HelpMethod == "" || len(input.TimeSlots) == 0 {
		return model.NeedDetail{}, errors.New(constants.ErrCodeInvalidPayload, constants.ErrMsgInvalidPayload)
	}
	entry, err := repository.SubmitResponse(needID, operator, input.HelpMethod, input.TimeSlots)
	if err != nil {
		return model.NeedDetail{}, err
	}
	logger.Info("响应加入候选名单", "need", needID, "responder", entry.Responder, "position", entry.Position)
	return GetNeedDetail(needID, operator)
}

// PickNeedCandidate 发布人从名单中挑一位约谈。
func PickNeedCandidate(needID int, input model.PickInput) (model.NeedDetail, error) {
	operator := resolveViewer(input.Operator)
	entry, err := repository.PickCandidate(needID, input.ResponseID, operator)
	if err != nil {
		return model.NeedDetail{}, err
	}
	logger.Info("安排约谈", "need", needID, "responder", entry.Responder, "position", entry.Position)
	return GetNeedDetail(needID, operator)
}

// ExitNeedResponse 响应者退出；若退出的是约谈对象，由下一位自动接手。
func ExitNeedResponse(needID int, input model.ExitInput) (model.NeedDetail, error) {
	operator := resolveViewer(input.Operator)
	entry, promoted, err := repository.ExitResponse(needID, operator)
	if err != nil {
		return model.NeedDetail{}, err
	}
	logger.Info("响应者离开队列", "need", needID, "responder", entry.Responder, "status", entry.Status)
	if promoted != nil {
		logger.Info("约谈自动接手", "need", needID, "responder", promoted.Responder, "position", promoted.Position)
	}
	return GetNeedDetail(needID, operator)
}

// CloseNeed 发布人关闭求助，名单停止变动。
func CloseNeed(needID int, input model.CloseInput) (model.NeedDetail, error) {
	operator := resolveViewer(input.Operator)
	if _, err := repository.CloseNeed(needID, operator); err != nil {
		return model.NeedDetail{}, err
	}
	logger.Info("求助已关闭，名单停止", "need", needID, "operator", operator)
	return GetNeedDetail(needID, operator)
}

// buildNeedSummary 汇总单条求助的名单状态与查看者视角信息。
func buildNeedSummary(need model.Need, queue []model.NeedResponse, viewer string) model.NeedSummary {
	summary := model.NeedSummary{
		Need:          need,
		StatusLabel:   constants.NeedStatusLabels[need.Status],
		ResponseCount: len(queue),
	}
	for i := range queue {
		entry := &queue[i]
		entry.StatusLabel = constants.ResponseStatusLabels[entry.Status]
		entry.Mine = entry.Responder == viewer
		switch entry.Status {
		case constants.ResponseStatusInterviewing:
			summary.CurrentInterview = entry.Responder
		case constants.ResponseStatusWaiting:
			summary.WaitingCount++
		}
	}
	summary.MyQueue = buildViewerQueue(queue, viewer)
	return summary
}

// buildViewerQueue 计算查看者自己的位次、前方人数与轮候结果。
func buildViewerQueue(queue []model.NeedResponse, viewer string) *model.ViewerQueue {
	var mine *model.NeedResponse
	for i := range queue {
		entry := &queue[i]
		if entry.Responder != viewer {
			continue
		}
		// 优先取仍在队列中的条目，否则取最近一条
		if mine == nil || isActiveStatus(entry.Status) || entry.Position > mine.Position {
			mine = entry
		}
	}
	if mine == nil {
		return nil
	}
	info := &model.ViewerQueue{
		ResponseID:  mine.ID,
		Position:    mine.Position,
		Status:      mine.Status,
		StatusLabel: constants.ResponseStatusLabels[mine.Status],
	}
	switch mine.Status {
	case constants.ResponseStatusWaiting:
		for i := range queue {
			if queue[i].Position < mine.Position && isActiveStatus(queue[i].Status) {
				info.AheadCount++
			}
		}
		info.Result = fmt.Sprintf(constants.QueueResultWaiting, info.AheadCount)
	case constants.ResponseStatusInterviewing:
		info.Result = constants.QueueResultInterviewing
	case constants.ResponseStatusWithdrawn:
		info.Result = constants.QueueResultWithdrawn
	case constants.ResponseStatusExited:
		info.Result = constants.QueueResultExited
	default:
		info.Result = constants.QueueResultClosed
	}
	return info
}

func isActiveStatus(status string) bool {
	return status == constants.ResponseStatusWaiting || status == constants.ResponseStatusInterviewing
}
