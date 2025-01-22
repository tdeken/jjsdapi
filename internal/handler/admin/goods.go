package admin

import (
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/fiber/result"
	meet "jjsdapi/internal/meet/admin"
)

// List 客户列表
// @Router /admin/goods/list [GET]
func (s Goods) List(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsListReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).List(form)
	return result.Json(ctx, res, err)
}

// Store 新增商品
// @Router /admin/goods/store [POST]
func (s Goods) Store(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsStoreReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Store(form)
	return result.Json(ctx, res, err)
}

// Update 更新商品
// @Router /admin/goods/update [POST]
func (s Goods) Update(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsUpdateReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Update(form)
	return result.Json(ctx, res, err)
}

// Destroy 删除商品
// @Router /admin/goods/destroy [POST]
func (s Goods) Destroy(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsDestroyReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).Destroy(form)
	return result.Json(ctx, res, err)
}
