package constants

// 求助状态：开放中可响应可约谈，关闭后候选名单停止
const (
	NeedStatusOpen   = "open"
	NeedStatusClosed = "closed"
)

// 候选人状态：排队中 -> 约谈中 -> 已退出
const (
	ResponseStatusWaiting      = "waiting"
	ResponseStatusInterviewing = "interviewing"
	ResponseStatusWithdrawn    = "withdrawn"
)

// 查看者轮候状态（未在名单中）
const (
	QueueStatusNone = "none"
)

// DateTimeLayout 候选名单提交时间的统一展示格式
const DateTimeLayout = "2006-01-02 15:04"
