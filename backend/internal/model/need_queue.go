package model

// NeedResponse 候选名单条目：响应者提交的能帮方式与可上门时段。
type NeedResponse struct {
	ID          int      `json:"id"`
	NeedID      int      `json:"needId"`
	Responder   string   `json:"responder"`
	HelpMethod  string   `json:"helpMethod"`
	TimeSlots   []string `json:"timeSlots"`
	Status      string   `json:"status"`
	StatusLabel string   `json:"statusLabel"`
	Position    int      `json:"position"` // 按提交先后的排队位次，从 1 开始
	SubmittedAt string   `json:"submittedAt"`
	Mine        bool     `json:"mine"` // 是否当前查看者本人提交
}

// ViewerQueue 查看者个人在名单中的位置与轮候结果。
type ViewerQueue struct {
	ResponseID  int    `json:"responseId"`
	Position    int    `json:"position"`
	Status      string `json:"status"`
	StatusLabel string `json:"statusLabel"`
	AheadCount  int    `json:"aheadCount"` // 前方仍在队列中的人数（含约谈中）
	Result      string `json:"result"`     // 轮候结果文案
}

// NeedSummary 求助列表视图：含当前约谈人与查看者自己的排队信息。
type NeedSummary struct {
	Need
	StatusLabel      string       `json:"statusLabel"`
	ResponseCount    int          `json:"responseCount"`
	WaitingCount     int          `json:"waitingCount"`
	CurrentInterview string       `json:"currentInterview"` // 当前约谈人，空串表示暂无
	MyQueue          *ViewerQueue `json:"myQueue"`
}

// NeedDetail 求助详情视图：完整候选名单与当前查看者可执行的操作。
type NeedDetail struct {
	NeedSummary
	Queue      []NeedResponse `json:"queue"`
	CanRespond bool           `json:"canRespond"` // 可提交响应
	CanPick    bool           `json:"canPick"`    // 发布人可挑人约谈
	CanClose   bool           `json:"canClose"`   // 发布人可关闭求助
}

// SubmitResponseInput 响应者提交内容。
type SubmitResponseInput struct {
	Operator   string   `json:"operator"`
	HelpMethod string   `json:"helpMethod"`
	TimeSlots  []string `json:"timeSlots"`
}

// PickInput 发布人挑选约谈对象。
type PickInput struct {
	Operator   string `json:"operator"`
	ResponseID int    `json:"responseId"`
}

// ExitInput 响应者退出/撤回。
type ExitInput struct {
	Operator string `json:"operator"`
}

// CloseInput 发布人关闭求助。
type CloseInput struct {
	Operator string `json:"operator"`
}
