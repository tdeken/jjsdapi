package service

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/app"
	"jjsdapi/internal/code"
	"jjsdapi/internal/fiber/mid"
	"jjsdapi/internal/repository/dao"
)

type Service struct {
	Ctx context.Context
	dao *dao.Query
	app *fiber.Ctx

	adminUser *mid.AdminAuth
}

func (s *Service) Init(ctx *fiber.Ctx) {
	s.app = ctx
	s.Ctx = ctx.Context()
	s.dao = dao.Q
}

func (s *Service) FiberCtx() *fiber.Ctx {
	return s.app
}

func (s *Service) GetDao() *dao.Query {
	return s.dao
}

func (s *Service) PushErr(err error, prr ...*code.Error) *code.Error {
	if err == nil {
		return nil
	}

	er, ok := code.As(err)
	if ok {
		return er
	}

	app.Logger.Error(s.Ctx, fmt.Sprintf("错误异常:%v", err))

	if len(prr) == 0 {
		return code.SystemError
	}

	return prr[0]
}

func (s *Service) ParamErr(tip string) *code.Error {
	return code.NewError(code.VerifyErrorCode, fmt.Sprintf("参数错误：%s", tip))
}

// AdminUser 获取用户信息
func (s *Service) AdminUser() *mid.AdminAuth {
	if s.adminUser == nil {
		user, ok := s.app.Context().UserValue(mid.AuthUser).(mid.AdminAuth)
		if ok {
			s.adminUser = &user
		}
	}

	if s.adminUser != nil {
		return s.adminUser
	}

	return &mid.AdminAuth{}
}

// AdminUserId 获取用户id
func (s *Service) AdminUserId() int64 {
	return s.AdminUser().UserId
}

// AdminTokenId 获取token ID
func (s *Service) AdminTokenId() string {
	return s.AdminUser().TokenId
}

// AdminToken 登陆token
func (s *Service) AdminToken() string {
	return s.AdminUser().Token
}
