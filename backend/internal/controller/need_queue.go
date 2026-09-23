package controller

import (
	"net/http"
	"strconv"

	bizerrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// NeedList 求助列表，viewer 用于计算“我的位置”和“当前约谈人”
func NeedList(c *gin.Context) {
	c.JSON(http.StatusOK, service.NeedSummaries(c.Query("viewer")))
}

// NeedDetail 求助详情：候选名单、当前约谈人和查看者的轮候结果
func NeedDetail(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	detail, err := service.NeedDetail(needID, c.Query("viewer"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

// SubmitNeedResponse 响应者提交候选：能帮的方式 + 可上门时段
func SubmitNeedResponse(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	var input model.SubmitResponseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, bizerrors.ErrInvalidResponseInput)
		return
	}
	detail, err := service.SubmitNeedResponse(needID, input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, detail)
}

// PickInterview 发布人挑一位排队中的响应者约谈
func PickInterview(c *gin.Context) {
	needID, responseID, ok := parseNeedAndResponseID(c)
	if !ok {
		return
	}
	operator, ok := bindOperator(c)
	if !ok {
		return
	}
	detail, err := service.PickInterview(needID, responseID, operator)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

// WithdrawNeedResponse 响应者退出轮候；约谈对象退出后由下一位自动接手
func WithdrawNeedResponse(c *gin.Context) {
	needID, responseID, ok := parseNeedAndResponseID(c)
	if !ok {
		return
	}
	operator, ok := bindOperator(c)
	if !ok {
		return
	}
	detail, err := service.WithdrawNeedResponse(needID, responseID, operator)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

// CloseNeed 发布人关闭求助，候选名单停止
func CloseNeed(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	operator, ok := bindOperator(c)
	if !ok {
		return
	}
	detail, err := service.CloseNeed(needID, operator)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

func parseNeedID(c *gin.Context) (int, bool) {
	needID, err := strconv.Atoi(c.Param("id"))
	if err != nil || needID <= 0 {
		respondError(c, bizerrors.ErrInvalidNeedID)
		return 0, false
	}
	return needID, true
}

func parseNeedAndResponseID(c *gin.Context) (int, int, bool) {
	needID, ok := parseNeedID(c)
	if !ok {
		return 0, 0, false
	}
	responseID, err := strconv.Atoi(c.Param("responseId"))
	if err != nil || responseID <= 0 {
		respondError(c, bizerrors.ErrResponseNotFound)
		return 0, 0, false
	}
	return needID, responseID, true
}

func bindOperator(c *gin.Context) (string, bool) {
	var input model.OperatorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, bizerrors.ErrMissingOperator)
		return "", false
	}
	return input.Operator, true
}
