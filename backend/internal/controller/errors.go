package controller

import (
	"net/http"

	bizerrors "cyskillswap/internal/errors"
	"github.com/gin-gonic/gin"
)

// respondError 统一把业务异常映射为 HTTP 响应
func respondError(c *gin.Context, err error) {
	if bizErr, ok := err.(bizerrors.BusinessError); ok {
		c.JSON(bizErr.Status, bizErr)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "服务器内部错误"})
}
