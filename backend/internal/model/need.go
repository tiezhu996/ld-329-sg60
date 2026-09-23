package model

// Need 求助信息，Status 为 closed 时候选名单停止
type Need struct {
	ID          int    `json:"id"`
	Requester   string `json:"requester"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Campus      string `json:"campus"`
	ExpectTime  string `json:"expectTime"`
	BudgetType  string `json:"budgetType"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Responses   int    `json:"responses"`
}

// NeedResponse 响应者提交的候选记录，按 Seq（提交先后）排队
type NeedResponse struct {
	ID          int      `json:"id"`
	NeedID      int      `json:"needId"`
	Responder   string   `json:"responder"`
	HelpOffer   string   `json:"helpOffer"`
	VisitSlots  []string `json:"visitSlots"`
	Status      string   `json:"status"`
	Seq         int      `json:"seq"`
	SubmittedAt string   `json:"submittedAt"`
}

// NeedSummary 列表行：在求助信息上附带当前约谈人和查看者自己的轮候结果
type NeedSummary struct {
	Need
	CurrentInterviewee string `json:"currentInterviewee"`
	WaitingCount       int    `json:"waitingCount"`
	IsPublisher        bool   `json:"isPublisher"`
	MyStatus           string `json:"myStatus"`
	MyPosition         int    `json:"myPosition"`
}

// NeedQueueEntry 候选名单条目，Position 为排队位次（约谈中/已退出为 0）
type NeedQueueEntry struct {
	ResponseID  int      `json:"responseId"`
	Position    int      `json:"position"`
	Responder   string   `json:"responder"`
	HelpOffer   string   `json:"helpOffer"`
	VisitSlots  []string `json:"visitSlots"`
	Status      string   `json:"status"`
	SubmittedAt string   `json:"submittedAt"`
	IsViewer    bool     `json:"isViewer"`
}

// NeedDetail 求助详情：摘要 + 完整候选名单
type NeedDetail struct {
	NeedSummary
	Entries []NeedQueueEntry `json:"entries"`
}

// SubmitResponseInput 响应求助的入参：能帮的方式 + 可上门时段
type SubmitResponseInput struct {
	Responder  string   `json:"responder"`
	HelpOffer  string   `json:"helpOffer"`
	VisitSlots []string `json:"visitSlots"`
}

// OperatorInput 需要操作人的请求入参（约谈/退出/关闭）
type OperatorInput struct {
	Operator string `json:"operator"`
}
