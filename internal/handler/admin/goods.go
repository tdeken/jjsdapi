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

// SelectInfo 表单要选择的数据
// @Router /admin/goods/select-info [POST]
func (s Goods) SelectInfo(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsSelectInfoReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).SelectInfo(form)
	return result.Json(ctx, res, err)
}

// SkuStore 创建销售品
// @Router /admin/goods/sku-store [POST]
func (s Goods) SkuStore(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsSkuStoreReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).SkuStore(form)
	return result.Json(ctx, res, err)
}

// SkuUpdate 更新销售品
// @Router /admin/goods/sku-update [POST]
func (s Goods) SkuUpdate(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsSkuUpdateReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).SkuUpdate(form)
	return result.Json(ctx, res, err)
}

// SkuDestroy 删除销售品
// @Router /admin/goods/sku-destroy [POST]
func (s Goods) SkuDestroy(ctx *fiber.Ctx) (e error) {
	var form = &meet.GoodsSkuDestroyReq{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}

	res, err := s.getDep(ctx).SkuDestroy(form)
	return result.Json(ctx, res, err)
}
