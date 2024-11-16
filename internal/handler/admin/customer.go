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
