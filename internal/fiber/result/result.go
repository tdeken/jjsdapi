package result

import (
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/code"
)

type WebResult struct {
	Code int32  `json:"code"` //错误码
	Msg  string `json:"msg"`  //返回的消息
	Data any    `json:"data"` //返回的数据结果
}

func Json(ctx *fiber.Ctx, res any, err error) error {
	var crr = code.Ok
	if err != nil {
		var ok bool
		crr, ok = code.As(err)
		if !ok {
			crr = code.NewError(code.CommonErrorCode, err.Error())
		}
	}

	var data = res
	if data == nil {
		data = make(map[string]interface{})
	}

	return ctx.JSON(WebResult{
		Code: crr.GetCode(),
		Msg:  crr.GetDetail(),
		Data: data,
	})
}
