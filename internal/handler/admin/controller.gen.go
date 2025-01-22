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
		Goods{},
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
		action.NewAction("POST", s.Logout, action.UseMidType("admin_jwt")),
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
		action.NewAction("GET", s.List, action.UseMidType("admin_jwt")),
		action.NewAction("GET", s.Select, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Store, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Update, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Destroy, action.UseMidType("admin_jwt")),
		action.NewAction("GET", s.AddressList, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.AddressCreate, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.AddressUpdate, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.AddressDestroy, action.UseMidType("admin_jwt")),
	}
}

// 获取依赖服务
func (s Customer) getDep(ctx *fiber.Ctx) admin.Customer {
	dep := admin.Customer{}
	dep.Init(ctx)
	return dep
}

// Goods 商品数据
type Goods struct {
	Controller
}

// Group 基础请求组
func (s Goods) Group() string {
	return "goods"
}

// Register 注册路由
func (s Goods) Register() []action.Action {
	return []action.Action{
		action.NewAction("GET", s.List, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Store, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Update, action.UseMidType("admin_jwt")),
		action.NewAction("POST", s.Destroy, action.UseMidType("admin_jwt")),
	}
}

// 获取依赖服务
func (s Goods) getDep(ctx *fiber.Ctx) admin.Goods {
	dep := admin.Goods{}
	dep.Init(ctx)
	return dep
}
