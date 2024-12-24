package admin

import (
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/fiber/result"
	meet "jjsdapi/internal/meet/admin"
)

// AddressList 配送地址列表
// @Router /admin/customer/address-list [GET]
func (s Customer) AddressList(ctx *fiber.Ctx) (e error) {
	var form = &meet.CustomerAddressListReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).AddressList(form)
	return result.Json(ctx, res, err)
}

// List 客户列表
// @Router /admin/customer/list [GET]
func (s Customer) List(ctx *fiber.Ctx) (e error) {
	var form = &meet.CustomerListReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).List(form)
	return result.Json(ctx, res, err)
}

// Select 客户列表选择
// @Router /admin/customer/select [GET]
func (s Customer) Select(ctx *fiber.Ctx) (e error) {
	var form = &meet.CustomerSelectReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Select(form)
	return result.Json(ctx, res, err)
}

// AddressCreate 地址创建
// @Router /admin/customer/address-create [POST]
func (s Customer) AddressCreate(ctx *fiber.Ctx) (e error) {
	var form = &meet.CustomerAddressCreateReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).AddressCreate(form)
	return result.Json(ctx, res, err)
}
