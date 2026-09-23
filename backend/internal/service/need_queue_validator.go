package service

import (
	"strings"

	bizerrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// validateResponseInput 校验响应入参：必须填写能帮的方式并至少选择一个可上门时段
func validateResponseInput(input model.SubmitResponseInput) error {
	if strings.TrimSpace(input.Responder) == "" {
		return bizerrors.ErrMissingResponder
	}
	if strings.TrimSpace(input.HelpOffer) == "" || len(input.VisitSlots) == 0 {
		return bizerrors.ErrInvalidResponseInput
	}
	return nil
}

// validateOperator 校验操作人（约谈/退出/关闭均需要）
func validateOperator(operator string) error {
	if strings.TrimSpace(operator) == "" {
		return bizerrors.ErrMissingOperator
	}
	return nil
}
