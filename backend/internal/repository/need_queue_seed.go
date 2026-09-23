package repository

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// seed 初始化求助与候选名单样例数据，排队顺序由提交先后（Seq）决定
func (s *NeedQueueStore) seed() {
	s.needs = []model.Need{
		{ID: 1, Requester: "孟野", Title: "找人帮忙拍乐队宣传照", Category: "摄影", Campus: "西校区", ExpectTime: "本周六上午", BudgetType: "技能交换", Description: "可交换 3 次吉他课，希望会调色和室外构图。", Status: constants.NeedStatusOpen},
		{ID: 2, Requester: "许安", Title: "求教 Python 数据分析", Category: "编程", Campus: "中心校区", ExpectTime: "周二晚", BudgetType: "小额报酬", Description: "论文问卷数据需要清洗和画图，最好有 pandas 经验。", Status: constants.NeedStatusOpen},
		{ID: 3, Requester: "林澈", Title: "想学吉他扫弦入门", Category: "乐器", Campus: "东校区", ExpectTime: "周三晚", BudgetType: "技能交换", Description: "用摄影课交换吉他基础，希望同校区或线上。", Status: constants.NeedStatusOpen},
	}

	s.addSeedResponse(1, "周芮", "可带全画幅相机和反光板上门，包调色修 9 图。", []string{"周六上午", "周六下午"}, constants.ResponseStatusInterviewing, "2026-09-20 10:24")
	s.addSeedResponse(1, "林澈", "提供毕业照同款人像拍摄，可交换吉他入门课。", []string{"周六上午"}, constants.ResponseStatusWaiting, "2026-09-20 15:41")
	s.addSeedResponse(1, "许安", "会基础构图和 Lightroom 调色，想积累乐队人像作品。", []string{"周六下午", "周日全天"}, constants.ResponseStatusWaiting, "2026-09-21 09:02")

	s.addSeedResponse(2, "周芮", "pandas 清洗加可视化一次搞定，可提供论文级代码模板。", []string{"周二晚"}, constants.ResponseStatusInterviewing, "2026-09-20 11:15")
	s.addSeedResponse(2, "孟野", "会 matplotlib 和问卷统计，可帮忙跑数据。", []string{"周二晚", "周三晚"}, constants.ResponseStatusWaiting, "2026-09-21 20:33")
	s.addSeedResponse(2, "林澈", "可以帮忙整理问卷数据，时间合适可上门。", []string{"周二晚"}, constants.ResponseStatusWithdrawn, "2026-09-22 08:47")

	s.addSeedResponse(3, "孟野", "民谣扫弦四课时入门，可交换人像摄影课。", []string{"周三晚", "周六上午"}, constants.ResponseStatusInterviewing, "2026-09-21 14:05")
	s.addSeedResponse(3, "周芮", "会基础弹唱，可以一起练习互相纠错。", []string{"周三晚"}, constants.ResponseStatusWaiting, "2026-09-22 19:18")
}

func (s *NeedQueueStore) addSeedResponse(needID int, responder, helpOffer string, visitSlots []string, status, submittedAt string) {
	s.responses = append(s.responses, model.NeedResponse{
		ID:          s.nextResponseID,
		NeedID:      needID,
		Responder:   responder,
		HelpOffer:   helpOffer,
		VisitSlots:  visitSlots,
		Status:      status,
		Seq:         s.nextSeq,
		SubmittedAt: submittedAt,
	})
	s.nextResponseID++
	s.nextSeq++
}
