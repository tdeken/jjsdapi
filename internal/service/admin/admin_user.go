package admin

import (
	"fmt"
	"jjsdapi/internal/app"
	"jjsdapi/internal/code"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/repository/ckey"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils/dbcheck"
	"jjsdapi/plugins/certs"
	"time"
)

type AdminUser struct {
	service.Service
}

// Login 登陆
func (s AdminUser) Login(req *meet.AdminUserLoginReq) (*meet.AdminUserLoginRes, error) {
	m := s.GetDao().AdminUser

	fmt.Println(s.Ctx)

	user, err := m.WithContext(s.Ctx).Where(m.Username.Eq(req.Username), m.DeletedAt.Eq(0)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if user == nil || user.Password != req.Password {
		return nil, code.AdminUserLoginErr
	}

	auth := certs.NewUser()
	token, err := auth.Token(&certs.User{
		UserId: user.ID,
	})
	if err != nil {
		return nil, s.PushErr(err)
	}

	exp := auth.Expired().Sub(time.Now())
	key := ckey.UserLogin(auth.UserId, auth.ID)
	err = app.Redis.Set(s.Ctx, key, token, exp)
	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.AdminUserLoginRes{
		Token: token,
		Name:  user.Name,
	}, nil
}

// Logout 登出
func (s AdminUser) Logout(req *meet.AdminUserLogoutReq) (*meet.AdminUserLogoutRes, error) {
	key := ckey.UserLogin(s.AdminUserId(), s.AdminTokenId())
	_ = app.Redis.Del(s.Ctx, key)
	return &meet.AdminUserLogoutRes{}, nil
}