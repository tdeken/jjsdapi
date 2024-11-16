package ip

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetRealIp(ctx *fiber.Ctx) (ip string) {
	ips := ctx.IPs()

	if len(ips) > 0 {
		return ips[0]
	}

	// 尝试从 X-Real-Ip 中获取
	ip = strings.TrimSpace(string(ctx.Request().Header.Peek(`X-Real-Ip`)))
	if ip == "" {
		return ctx.IP()
	}

	return
}
