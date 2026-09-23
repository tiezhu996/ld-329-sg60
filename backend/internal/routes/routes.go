package routes

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/controller"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group(constants.APIPrefix)
	api.GET("/health", controller.Health)
	api.GET("/dashboard/overview", controller.Overview)
	api.GET("/skills", controller.Skills)

	// 技能求助与候选名单
	api.GET("/needs", controller.NeedSummaries)
	api.GET("/needs/:id", controller.NeedDetail)
	api.POST("/needs/:id/responses", controller.SubmitNeedResponse)
	api.POST("/needs/:id/pick", controller.PickNeedCandidate)
	api.POST("/needs/:id/exit", controller.ExitNeedResponse)
	api.POST("/needs/:id/close", controller.CloseNeed)

	api.GET("/matches", controller.Matches)
	api.GET("/appointments", controller.Appointments)
	api.GET("/reviews", controller.Reviews)
	api.GET("/messages", controller.Messages)
	api.GET("/profile", controller.Profile)
}
