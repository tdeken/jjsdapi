package mid

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/xid"
	"jjsdapi/internal/app"
	"jjsdapi/internal/code"
	"jjsdapi/internal/config"
	"jjsdapi/internal/consts"
	"jjsdapi/internal/fiber/ip"
	"jjsdapi/internal/fiber/result"
	"time"
)

func Recover() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// 捕获后发送到微信处理
				var errMsg string
				errMsg += "环境:" + config.Conf.Server.Env + "\n"
				errMsg += "请求服务:[" + string(ctx.Request().Header.Method()) + "]" + string(ctx.Request().RequestURI()) + "\n"
				rawData := ctx.Request().Body()
				errMsg += "请求rawData:" + fmt.Sprint(string(rawData)) + "\n"
				errMsg += "错误内容:" + fmt.Sprintf("%s", r)
				app.Logger.ErrorWithStack(ctx.Context(), errMsg)
				err = nil
				_ = ctx.JSON(result.WebResult{
					Code: code.SystemErrorCode,
					Msg:  code.SystemError.GetDetail(),
					Data: nil,
				})
			}
		}()

		return ctx.Next()
	}
}

func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		t := time.Now()
		defer func() {
			// 请求后
			latency := time.Since(t)
			status := ctx.Response().StatusCode()
			app.Logger.Info(ctx.Context(), fmt.Sprintf("[Request] %v | %3d | %12v | %s | %-7s %s, res: %s",
				t.Format(time.DateTime),
				status,
				latency,
				ip.GetRealIp(ctx),
				ctx.Request().Header.Method(),
				ctx.Request().URI().Path(),
				string(ctx.Response().Body()),
			))
		}()

		return ctx.Next()
	}
}

func TraceId() fiber.Handler {
	return requestid.New(
		requestid.Config{
			Generator: func() string {
				return xid.New().String()
			},
			ContextKey: consts.TraceIdKey,
		},
	)
}

func Cors() fiber.Handler {
	return cors.New()
}
