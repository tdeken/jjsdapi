package handler

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	action "github.com/tdeken/fiberaction"
	"jjsdapi/internal/code"
	"jjsdapi/internal/fiber/mid"
	"jjsdapi/internal/fiber/validate"
	"strings"
)

type Handler struct {
}

// ValidateRequest 统一校验请求数据
func (s Handler) ValidateRequest(ctx *fiber.Ctx, rt validate.RequestInterface) *code.Error {
	err := validate.CheckParams(ctx, rt)
	if err != nil {
		var errMsg string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, er := range err.(validator.ValidationErrors) {
				errMsg = fmt.Sprintf("错误字段:%v,验证类型:%v:%v,参数值:%v", er.Field(), er.Tag(), er.Param(), er.Value())
				break
			}
		} else {
			errMsg = err.Error()
		}
		return code.NewError(code.VerifyErrorCode, errMsg)
	}

	return nil
}

// ChooseMid 可以选择的服务中间件
func (s Handler) ChooseMid(t action.MidType) (ms []fiber.Handler) {
	if t == nil {
		return
	}

	mst := fmt.Sprint(t)
	msa := strings.Split(mst, ",")

	for _, v := range msa {
		if m, ok := mid.Map[v]; ok {
			ms = append(ms, m)
		}

	}

	return
}
