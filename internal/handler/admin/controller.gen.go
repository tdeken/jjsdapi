//DO NOT EDIT.

package admin

import (
	"github.com/gofiber/fiber/v2"
	action "github.com/tdeken/fiberaction"
	"jjsdapi/internal/fiber/server"
	"jjsdapi/internal/service/admin"
)

// Route 模块路由
func Route() {
	r := server.Web.Server.Group("admin")

	action.AutoRegister(r,
		AdminUser{},
		Customer{},
	)

}

// AdminUser 后台用户
type AdminUser struct {
	Controller
}

// Group 基础请求组
func (s AdminUser) Group() string {
	return "admin-user"
}

// Register 注册路由
func (s AdminUser) Register() []action.Action {
	return []action.Action{
		action.NewAction("POST", s.Login),
		action.NewAction("GET", s.Logout, action.UseMidType("admin_jwt")),
		action.NewAction("GET", s.Info, action.UseMidType("admin_jwt")),
	}
}

// 获取依赖服务
func (s AdminUser) getDep(ctx *fiber.Ctx) admin.AdminUser {
	dep := admin.AdminUser{}
	dep.Init(ctx)
	return dep
}

// Customer 客户数据
type Customer struct {
	Controller
}

// Group 基础请求组
func (s Customer) Group() string {
	return "customer"
}

// Register 注册路由
func (s Customer) Register() []action.Action {
	return []action.Action{
		action.NewAction("GET", s.AddressList),
	}
}

// 获取依赖服务
func (s Customer) getDep(ctx *fiber.Ctx) admin.Customer {
	dep := admin.Customer{}
	dep.Init(ctx)
	return dep
}
