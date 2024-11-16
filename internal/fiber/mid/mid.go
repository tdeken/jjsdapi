package mid

import "github.com/gofiber/fiber/v2"

var Map = map[string]fiber.Handler{
	"admin_jwt": AdminJwt(),
}

const (
	headerAuthKey    = "Authorization"
	headerAuthScheme = "Bearer"
	AuthUser         = "__auth__user__"
)
