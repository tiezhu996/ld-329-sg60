package constants

// 求助状态
const (
	NeedStatusOpen   = "open"
	NeedStatusClosed = "closed"
)

// 求助状态展示文案
var NeedStatusLabels = map[string]string{
	NeedStatusOpen:   "开放中",
	NeedStatusClosed: "已关闭",
}

// 候选名单（响应）状态
const (
	ResponseStatusWaiting      = "waiting"      // 排队中
	ResponseStatusInterviewing = "interviewing" // 约谈中
	ResponseStatusWithdrawn    = "withdrawn"    // 约谈对象撤回
	ResponseStatusExited       = "exited"       // 约谈开始前主动退出
	ResponseStatusClosed       = "closed"       // 求助关闭，轮候结束
)

// 候选状态展示文案
var ResponseStatusLabels = map[string]string{
	ResponseStatusWaiting:      "排队中",
	ResponseStatusInterviewing: "约谈中",
	ResponseStatusWithdrawn:    "已撤回",
	ResponseStatusExited:       "已退出",
	ResponseStatusClosed:       "轮候结束",
}

// DefaultViewer 未指定用户时的默认视角（演示账号）
const DefaultViewer = "林澈"

// 轮候结果文案模板
const (
	QueueResultWaiting      = "排队中，前方还有 %d 位"
	QueueResultInterviewing = "约谈中，请等待发布人联系"
	QueueResultWithdrawn    = "已撤回约谈，名单由下一位接手"
	QueueResultExited       = "已退出排队"
	QueueResultClosed       = "求助已关闭，轮候结束"
	QueueResultNeedClosed   = "求助已关闭，停止响应"
)

// 业务错误码
const (
	ErrCodeNeedNotFound        = "NEED_NOT_FOUND"
	ErrCodeNeedClosed          = "NEED_CLOSED"
	ErrCodeSelfResponse        = "SELF_RESPONSE"
	ErrCodeDuplicateResponse   = "DUPLICATE_RESPONSE"
	ErrCodeNotRequester        = "NOT_REQUESTER"
	ErrCodeNotResponder        = "NOT_RESPONDER"
	ErrCodeResponseNotFound    = "RESPONSE_NOT_FOUND"
	ErrCodeResponseNotActive   = "RESPONSE_NOT_ACTIVE"
	ErrCodeCandidateNotWaiting = "CANDIDATE_NOT_WAITING"
	ErrCodeInterviewInProgress = "INTERVIEW_IN_PROGRESS"
	ErrCodeInvalidPayload      = "INVALID_PAYLOAD"
)

// 业务错误文案
const (
	ErrMsgNeedNotFound        = "求助不存在"
	ErrMsgNeedClosed          = "求助已关闭，名单停止变动"
	ErrMsgSelfResponse        = "不能响应自己发布的求助"
	ErrMsgDuplicateResponse   = "你已在候选名单中，请勿重复提交"
	ErrMsgNotRequester        = "只有发布人可以执行该操作"
	ErrMsgNotResponder        = "只能退出自己的响应"
	ErrMsgResponseNotFound    = "响应记录不存在"
	ErrMsgResponseNotActive   = "该响应已不在队列中"
	ErrMsgCandidateNotWaiting = "该候选人不在排队状态"
	ErrMsgInterviewInProgress = "当前已有约谈进行中，同一时间只能安排一位"
	ErrMsgInvalidPayload      = "提交内容不完整：请填写能帮的方式和可上门时段"
)
