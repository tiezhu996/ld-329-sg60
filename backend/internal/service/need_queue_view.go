package service

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// buildNeedSummary 汇总单条求助：当前约谈人、排队人数和查看者自己的轮候结果
func buildNeedSummary(need model.Need, responses []model.NeedResponse, viewer string) model.NeedSummary {
	summary := model.NeedSummary{Need: need, IsPublisher: viewer != "" && viewer == need.Requester}
	for _, resp := range responses {
		switch resp.Status {
		case constants.ResponseStatusInterviewing:
			summary.CurrentInterviewee = resp.Responder
		case constants.ResponseStatusWaiting:
			summary.WaitingCount++
		}
	}
	summary.MyStatus, summary.MyPosition = viewerQueueState(responses, viewer)
	return summary
}

// buildNeedDetail 组装详情：摘要 + 候选名单（活跃者按提交先后，已退出排在末尾）
func buildNeedDetail(need model.Need, responses []model.NeedResponse, viewer string) model.NeedDetail {
	return model.NeedDetail{
		NeedSummary: buildNeedSummary(need, responses, viewer),
		Entries:     buildQueueEntries(responses, viewer),
	}
}

// viewerQueueState 计算查看者在名单中的轮候结果：状态和排队位次
func viewerQueueState(responses []model.NeedResponse, viewer string) (string, int) {
	if viewer == "" {
		return constants.QueueStatusNone, 0
	}
	position := 0
	status := constants.QueueStatusNone
	for _, resp := range responses {
		if resp.Status == constants.ResponseStatusWaiting {
			position++
			if resp.Responder == viewer {
				status = constants.ResponseStatusWaiting
				return status, position
			}
			continue
		}
		if resp.Responder == viewer && resp.Status == constants.ResponseStatusInterviewing {
			return constants.ResponseStatusInterviewing, 0
		}
		if resp.Responder == viewer && resp.Status == constants.ResponseStatusWithdrawn {
			status = constants.ResponseStatusWithdrawn
		}
	}
	return status, 0
}

// buildQueueEntries 生成候选名单条目：约谈中与排队中按提交先后排列并标注位次，已退出沉底
func buildQueueEntries(responses []model.NeedResponse, viewer string) []model.NeedQueueEntry {
	active := make([]model.NeedResponse, 0, len(responses))
	withdrawn := make([]model.NeedResponse, 0, len(responses))
	for _, resp := range responses {
		if resp.Status == constants.ResponseStatusWithdrawn {
			withdrawn = append(withdrawn, resp)
			continue
		}
		active = append(active, resp)
	}
	entries := make([]model.NeedQueueEntry, 0, len(responses))
	position := 0
	for _, resp := range active {
		entry := toQueueEntry(resp, viewer)
		if resp.Status == constants.ResponseStatusWaiting {
			position++
			entry.Position = position
		}
		entries = append(entries, entry)
	}
	for _, resp := range withdrawn {
		entries = append(entries, toQueueEntry(resp, viewer))
	}
	return entries
}

func toQueueEntry(resp model.NeedResponse, viewer string) model.NeedQueueEntry {
	return model.NeedQueueEntry{
		ResponseID:  resp.ID,
		Responder:   resp.Responder,
		HelpOffer:   resp.HelpOffer,
		VisitSlots:  resp.VisitSlots,
		Status:      resp.Status,
		SubmittedAt: resp.SubmittedAt,
		IsViewer:    viewer != "" && resp.Responder == viewer,
	}
}
