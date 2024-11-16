package admin

import (
	"github.com/gofiber/fiber/v2"
	action "github.com/tdeken/fiberaction"
	"jjsdapi/internal/handler"
)

type Controller struct {
	handler.Handler
}

// ChooseMid 可以选择的服务中间件
func (c Controller) ChooseMid(t action.MidType) []fiber.Handler {
	switch t {
	default:
		return nil
	}
}
