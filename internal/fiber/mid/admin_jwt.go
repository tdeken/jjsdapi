package mid

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/app"
	"jjsdapi/internal/code"
	"jjsdapi/internal/fiber/result"
	"jjsdapi/internal/repository/ckey"
	"jjsdapi/plugins/certs"
	"strings"
)

type AdminAuth struct {
	UserId  int64  //员工id
	TokenId string //token ID
	Token   string //登陆的token
}

func AdminJwt() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get(headerAuthKey)
		l := len(headerAuthScheme)

		if len(auth) <= l+1 || !strings.EqualFold(auth[:l], headerAuthScheme) {
			return result.Json(ctx, nil, code.AuthError)
		}

		token := strings.TrimSpace(auth[l:])
		if token == "" {
			return result.Json(ctx, nil, code.AuthError)
		}

		user, fail := isCertAuthFail(ctx.Context(), token)
		if fail {
			return result.Json(ctx, nil, code.AuthError)
		}

		ctx.Context().SetUserValue(AuthUser, user)

		return ctx.Next()
	}
}

// 是否是鉴权失败
func isCertAuthFail(ctx context.Context, token string) (AdminAuth, bool) {
	auth := certs.NewUser()
	err := auth.Parse(token)
	if err != nil {
		return AdminAuth{}, true
	}

	key := ckey.UserLogin(auth.UserId, auth.ID)
	exist, err := app.Redis.Exists(ctx, key)
	if err != nil || !exist {
		return AdminAuth{}, true
	}

	return AdminAuth{
		UserId:  auth.UserId,
		TokenId: auth.ID,
		Token:   token,
	}, false
}
