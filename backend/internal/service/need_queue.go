package service

import (
	bizerrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

// NeedSummaries 求助列表：每条求助附带当前约谈人和查看者自己的轮候结果
func NeedSummaries(viewer string) []model.NeedSummary {
	needs := repository.ListNeeds()
	summaries := make([]model.NeedSummary, 0, len(needs))
	for _, need := range needs {
		_, responses, _ := repository.NeedSnapshot(need.ID)
		summaries = append(summaries, buildNeedSummary(need, responses, viewer))
	}
	return summaries
}

// NeedDetail 求助详情：摘要 + 按提交先后排队的候选名单
func NeedDetail(needID int, viewer string) (model.NeedDetail, error) {
	need, responses, ok := repository.NeedSnapshot(needID)
	if !ok {
		return model.NeedDetail{}, bizerrors.ErrNeedNotFound
	}
	return buildNeedDetail(need, responses, viewer), nil
}

// SubmitNeedResponse 响应者提交候选（能帮的方式 + 可上门时段），按提交先后入队
func SubmitNeedResponse(needID int, input model.SubmitResponseInput) (model.NeedDetail, error) {
	if err := validateResponseInput(input); err != nil {
		return model.NeedDetail{}, err
	}
	if _, err := repository.SubmitResponse(needID, input); err != nil {
		return model.NeedDetail{}, err
	}
	return NeedDetail(needID, input.Responder)
}

// PickInterview 发布人从候选名单中挑一位约谈，同一求助同时只安排一个约谈
func PickInterview(needID, responseID int, operator string) (model.NeedDetail, error) {
	if err := validateOperator(operator); err != nil {
		return model.NeedDetail{}, err
	}
	if err := repository.PickInterview(needID, responseID, operator); err != nil {
		return model.NeedDetail{}, err
	}
	return NeedDetail(needID, operator)
}

// WithdrawNeedResponse 响应者退出轮候；若退出的是当前约谈人，下一位自动接手
func WithdrawNeedResponse(needID, responseID int, operator string) (model.NeedDetail, error) {
	if err := validateOperator(operator); err != nil {
		return model.NeedDetail{}, err
	}
	if err := repository.WithdrawResponse(needID, responseID, operator); err != nil {
		return model.NeedDetail{}, err
	}
	return NeedDetail(needID, operator)
}

// CloseNeed 发布人关闭求助，候选名单停止
func CloseNeed(needID int, operator string) (model.NeedDetail, error) {
	if err := validateOperator(operator); err != nil {
		return model.NeedDetail{}, err
	}
	if err := repository.CloseNeed(needID, operator); err != nil {
		return model.NeedDetail{}, err
	}
	return NeedDetail(needID, operator)
}
