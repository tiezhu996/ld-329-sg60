package errors

import "net/http"

// 求助候选名单相关错误，错误码与提示信息集中维护
var (
	ErrNeedNotFound         = New(http.StatusNotFound, "NEED_NOT_FOUND", "求助不存在")
	ErrNeedClosed           = New(http.StatusConflict, "NEED_CLOSED", "求助已关闭，候选名单已停止")
	ErrNeedAlreadyClosed    = New(http.StatusConflict, "NEED_ALREADY_CLOSED", "求助已关闭，请勿重复操作")
	ErrNotPublisher         = New(http.StatusForbidden, "NOT_PUBLISHER", "只有发布人可以执行该操作")
	ErrResponseNotFound     = New(http.StatusNotFound, "RESPONSE_NOT_FOUND", "响应记录不存在")
	ErrDuplicateResponse    = New(http.StatusConflict, "DUPLICATE_RESPONSE", "你已在候选名单中，请勿重复响应")
	ErrRespondOwnNeed       = New(http.StatusConflict, "RESPOND_OWN_NEED", "不能响应自己发布的求助")
	ErrInterviewOccupied    = New(http.StatusConflict, "INTERVIEW_OCCUPIED", "当前已有约谈进行中，同一求助同时只能安排一个约谈")
	ErrCandidateNotWaiting  = New(http.StatusConflict, "CANDIDATE_NOT_WAITING", "该响应者不在排队状态，无法安排约谈")
	ErrNotCandidate         = New(http.StatusForbidden, "NOT_CANDIDATE", "只有响应者本人可以退出轮候")
	ErrAlreadyWithdrawn     = New(http.StatusConflict, "ALREADY_WITHDRAWN", "该响应已退出，请勿重复操作")
	ErrInvalidResponseInput = New(http.StatusBadRequest, "INVALID_RESPONSE_INPUT", "请填写能帮的方式并至少选择一个可上门时段")
	ErrMissingResponder     = New(http.StatusBadRequest, "MISSING_RESPONDER", "缺少响应者信息")
	ErrMissingOperator      = New(http.StatusBadRequest, "MISSING_OPERATOR", "缺少操作人信息")
	ErrInvalidNeedID        = New(http.StatusBadRequest, "INVALID_NEED_ID", "求助编号不合法")
)
