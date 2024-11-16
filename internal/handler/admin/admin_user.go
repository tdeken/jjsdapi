package admin

import (
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/fiber/result"
	meet "jjsdapi/internal/meet/admin"
)

// Login 登陆
// @Router /admin/admin-user/login [POST]
func (s AdminUser) Login(ctx *fiber.Ctx) (e error) {
	var form = &meet.AdminUserLoginReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Login(form)
	return result.Json(ctx, res, err)
}

// Logout 登出
// @Router /admin/admin-user/logout [GET]
func (s AdminUser) Logout(ctx *fiber.Ctx) (e error) {
	var form = &meet.AdminUserLogoutReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Logout(form)
	return result.Json(ctx, res, err)
}
