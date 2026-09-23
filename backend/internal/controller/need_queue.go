package controller

import (
	"net/http"
	"strconv"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// NeedSummaries 求助列表，携带查看者视角的排队信息。
func NeedSummaries(c *gin.Context) {
	c.JSON(http.StatusOK, service.ListNeedSummaries(c.Query("user")))
}

// NeedDetail 求助详情：完整候选名单。
func NeedDetail(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	detail, err := service.GetNeedDetail(needID, c.Query("user"))
	respondDetail(c, detail, err)
}

// SubmitNeedResponse 响应者提交能帮的方式与可上门时段。
func SubmitNeedResponse(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	var input model.SubmitResponseInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := service.SubmitNeedResponse(needID, input)
	respondDetail(c, detail, err)
}

// PickNeedCandidate 发布人从名单中挑一位约谈。
func PickNeedCandidate(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	var input model.PickInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := service.PickNeedCandidate(needID, input)
	respondDetail(c, detail, err)
}

// ExitNeedResponse 响应者退出；约谈对象退出记为撤回并由下一位接手。
func ExitNeedResponse(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	var input model.ExitInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := service.ExitNeedResponse(needID, input)
	respondDetail(c, detail, err)
}

// CloseNeed 发布人关闭求助，名单停止。
func CloseNeed(c *gin.Context) {
	needID, ok := parseNeedID(c)
	if !ok {
		return
	}
	var input model.CloseInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := service.CloseNeed(needID, input)
	respondDetail(c, detail, err)
}

func parseNeedID(c *gin.Context) (int, bool) {
	needID, err := strconv.Atoi(c.Param("id"))
	if err != nil || needID <= 0 {
		respondError(c, errors.New(constants.ErrCodeInvalidPayload, constants.ErrMsgInvalidPayload))
		return 0, false
	}
	return needID, true
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		respondError(c, errors.New(constants.ErrCodeInvalidPayload, constants.ErrMsgInvalidPayload))
		return false
	}
	return true
}

func respondDetail(c *gin.Context, detail model.NeedDetail, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

// respondError 业务异常按错误码映射状态码，未知异常返回 500。
func respondError(c *gin.Context, err error) {
	if biz, ok := err.(errors.BusinessError); ok {
		status := http.StatusBadRequest
		if biz.Code == constants.ErrCodeNeedNotFound || biz.Code == constants.ErrCodeResponseNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, biz)
		return
	}
	c.JSON(http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "服务内部错误"))
}
