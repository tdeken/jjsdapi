package admin

import (
	"fmt"
	"gorm.io/gen"
	"gorm.io/gen/field"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/repository/model"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils"
	"jjsdapi/internal/utils/dbcheck"
	"jjsdapi/internal/utils/timez"
	"time"
)

type Goods struct {
	service.Service
}

// List 客户列表
func (s Goods) List(req *meet.GoodsListReq) (*meet.GoodsListRes, error) {
	m := s.GetDao().Good

	var conds = []gen.Condition{
		m.DeletedAt.Eq(0),
	}

	if req.Keyword != "" {
		conds = append(conds, m.WithContext(s.Ctx).
			Or(m.Title.Like(utils.DbAllLike(req.Keyword))).
			Or(m.AsTitle.Like(utils.DbAllLike(req.Keyword))))
	}

	if req.Start != "" {
		conds = append(conds, m.UpdatedAt.Gte(timez.TableSearchTime(req.Start)))
	}

	if req.End != "" {
		conds = append(conds, m.UpdatedAt.Gte(timez.TableSearchTime(req.End)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PageSize)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.UpdatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, s.PushErr(err)
	}

	if total == 0 {
		return &meet.GoodsListRes{}, nil
	}

	var ids = make([]int64, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.ID)
	}

	skuM := s.GetDao().GoodsSku
	skuList, err := skuM.WithContext(s.Ctx).Where(skuM.GoodsID.In(ids...), skuM.DeletedAt.Eq(0)).Order(skuM.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, s.PushErr(err)
	}

	var skuMap = make(map[int64][]*model.GoodsSku)
	for _, v := range skuList {
		skuMap[v.GoodsID] = append(skuMap[v.GoodsID], v)
	}

	var data = make([]*meet.GoodsListOne, 0, len(list))
	for _, v := range list {
		one := &meet.GoodsListOne{
			Id:          utils.LongNumIdToStr(v.ID),
			Title:       v.Title,
			AsTitle:     v.AsTitle,
			SkuNum:      int64(len(skuMap[v.ID])),
			UpdatedDate: timez.TableDateTime(v.UpdatedAt),
			GoodsSkus:   make([]*meet.GoodsListOneGoodsSku, 0, len(skuMap[v.ID])),
		}
		for _, sku := range skuMap[v.ID] {
			name := v.Title
			if sku.Capacity != "" {
				name = sku.Capacity + name
			}

			if sku.Remark != "" {
				name += fmt.Sprintf("(%s)", sku.Remark)
			}

			one.GoodsSkus = append(one.GoodsSkus, &meet.GoodsListOneGoodsSku{
				Id:       utils.LongNumIdToStr(sku.ID),
				GoodsId:  one.Id,
				Name:     name,
				Capacity: sku.Capacity,
				Remark:   sku.Remark,
				Format:   sku.Format,
				Unit:     sku.Unit,
				Pp:       utils.Price(sku.Pp),
				Wp:       utils.Price(sku.Wp),
				Rp:       utils.Price(sku.Rp),
				Stock:    sku.Stock,
				Number:   sku.Number,
			})
		}

		data = append(data, one)
	}

	return &meet.GoodsListRes{
		List:  data,
		Total: total,
	}, nil
}

// Store 新增商品
func (s Goods) Store(req *meet.GoodsStoreReq) (*meet.GoodsStoreRes, error) {
	m := s.GetDao().Good
	exist, err := m.WithContext(s.Ctx).Where(m.Title.Eq(req.Title), m.DeletedAt.Eq(0)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if exist != nil {
		return nil, s.ParamErr("商品已存在")
	}

	err = m.WithContext(s.Ctx).Create(&model.Good{
		Title:   req.Title,
		AsTitle: req.AsTitle,
	})
	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.GoodsStoreRes{}, nil
}

// Update 更新商品
func (s Goods) Update(req *meet.GoodsUpdateReq) (*meet.GoodsUpdateRes, error) {
	id := utils.StrToLongNumId(req.Id)

	m := s.GetDao().Good
	exist, err := m.WithContext(s.Ctx).Where(m.Title.Eq(req.Title), m.DeletedAt.Eq(0), m.ID.Neq(id)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if exist != nil {
		return nil, s.ParamErr("商品已存在")
	}

	goods, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if goods == nil {
		return nil, s.ParamErr("商品不存在")
	}

	var up []field.AssignExpr
	if goods.Title != req.Title {
		up = append(up, m.Title.Value(req.Title))
	}
	if goods.AsTitle != req.AsTitle {
		up = append(up, m.AsTitle.Value(req.AsTitle))
	}

	if len(up) > 0 {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(goods.ID)).UpdateSimple(up...)
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.GoodsUpdateRes{}, nil
}

// Destroy 删除商品
func (s Goods) Destroy(req *meet.GoodsDestroyReq) (*meet.GoodsDestroyRes, error) {
	id := utils.StrToLongNumId(req.Id)
	m := s.GetDao().Good
	goods, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if goods != nil && goods.DeletedAt == 0 {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(id)).UpdateSimple(m.DeletedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.GoodsDestroyRes{}, nil
}

// SelectInfo 表单要选择的数据
func (s Goods) SelectInfo(req *meet.GoodsSelectInfoReq) (*meet.GoodsSelectInfoRes, error) {
	var format = []*meet.GoodsSelectInfoOne{
		{Label: "1*1", Value: "1*1"},
		{Label: "1*2", Value: "1*2"},
		{Label: "1*4", Value: "1*4"},
		{Label: "1*6", Value: "1*6"},
		{Label: "1*8", Value: "1*8"},
		{Label: "1*10", Value: "1*10"},
		{Label: "1*12", Value: "1*12"},
		{Label: "1*15", Value: "1*15"},
		{Label: "1*20", Value: "1*20"},
		{Label: "1*24", Value: "1*24"},
		{Label: "1*16", Value: "1*16"},
		{Label: "1*40", Value: "1*40"},
		{Label: "1*36", Value: "1*36"},
		{Label: "1*2*20", Value: "1*2*20"},
		{Label: "1*100", Value: "1*100"},
		{Label: "1*50", Value: "1*50"},
		{Label: "1*32", Value: "1*32"},
		{Label: "1*30", Value: "1*30"},
		{Label: "1*60", Value: "1*60"},
		{Label: "1*80", Value: "1*80"},
	}

	var unit = []*meet.GoodsSelectInfoOne{
		{Label: "件", Value: "件"},
		{Label: "份", Value: "份"},
		{Label: "箱", Value: "箱"},
		{Label: "斤", Value: "斤"},
		{Label: "瓶", Value: "瓶"},
		{Label: "条", Value: "条"},
		{Label: "包", Value: "包"},
		{Label: "个", Value: "个"},
		{Label: "袋", Value: "袋"},
		{Label: "盒", Value: "盒"},
		{Label: "捆", Value: "捆"},
		{Label: "抽", Value: "抽"},
		{Label: "提", Value: "提"},
		{Label: "支", Value: "支"},
		{Label: "捆", Value: "捆"},
	}

	return &meet.GoodsSelectInfoRes{
		Format: format,
		Unit:   unit,
	}, nil
}

// SkuStore 创建销售品
func (s Goods) SkuStore(req *meet.GoodsSkuStoreReq) (*meet.GoodsSkuStoreRes, error) {
	goodsId := utils.StrToLongNumId(req.GoodsId)

	val := fmt.Sprintf("%d_%s_%s_%s_%s", goodsId, req.Capacity, req.Remark, req.Format, req.Unit)
	mark, err := utils.Md5Str([]byte(val))
	if err != nil {
		return nil, s.PushErr(err)
	}

	m := s.GetDao().GoodsSku
	exist, err := m.WithContext(s.Ctx).Where(m.Mark.Eq(mark), m.DeletedAt.Eq(0)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}
	if exist != nil {
		return nil, s.ParamErr("销售品已存在")
	}

	err = m.WithContext(s.Ctx).Create(&model.GoodsSku{
		GoodsID:  goodsId,
		Mark:     mark,
		Capacity: req.Capacity,
		Remark:   req.Remark,
		Format:   req.Format,
		Unit:     req.Unit,
		Pp:       utils.PriceNumber(req.Pp),
		Wp:       utils.PriceNumber(req.Wp),
		Rp:       utils.PriceNumber(req.Rp),
		Stock:    req.Stock,
		Number:   req.Number,
	})

	if err != nil {
		return nil, s.PushErr(err)
	}

	goodsM := s.GetDao().Good
	_, err = goodsM.WithContext(s.Ctx).Where(goodsM.ID.Eq(goodsId)).UpdateSimple(goodsM.UpdatedAt.Value(time.Now().Unix()))
	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.GoodsSkuStoreRes{}, nil
}

// SkuUpdate 更新销售品
func (s Goods) SkuUpdate(req *meet.GoodsSkuUpdateReq) (*meet.GoodsSkuUpdateRes, error) {
	id := utils.StrToLongNumId(req.Id)
	m := s.GetDao().GoodsSku
	sku, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}
	if sku == nil {
		return nil, s.ParamErr("销售品不存在")
	}

	val := fmt.Sprintf("%d_%s_%s_%s_%s", sku.GoodsID, req.Capacity, req.Remark, req.Format, req.Unit)
	mark, err := utils.Md5Str([]byte(val))
	if err != nil {
		return nil, s.PushErr(err)
	}

	exist, err := m.WithContext(s.Ctx).Where(m.Mark.Eq(mark), m.ID.Neq(sku.ID), m.DeletedAt.Eq(0)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}
	if exist != nil {
		return nil, s.ParamErr("销售品已存在")
	}

	pp := utils.PriceNumber(req.Pp)
	wp := utils.PriceNumber(req.Wp)
	rp := utils.PriceNumber(req.Rp)

	var up []field.AssignExpr

	if sku.Mark != mark {
		up = append(up, m.Mark.Value(mark))
	}

	if sku.Capacity != req.Capacity {
		up = append(up, m.Capacity.Value(req.Capacity))
	}

	if sku.Remark != req.Remark {
		up = append(up, m.Remark.Value(req.Remark))
	}

	if sku.Format != req.Format {
		up = append(up, m.Format.Value(req.Format))
	}

	if sku.Unit != req.Unit {
		up = append(up, m.Unit.Value(req.Unit))
	}

	if sku.Pp != pp {
		up = append(up, m.Pp.Value(pp))
	}

	if sku.Wp != wp {
		up = append(up, m.Wp.Value(wp))
	}

	if sku.Rp != rp {
		up = append(up, m.Rp.Value(rp))
	}

	if sku.Stock != req.Stock {
		up = append(up, m.Stock.Value(req.Stock))
	}

	if sku.Number != req.Number {
		up = append(up, m.Number.Value(req.Number))
	}

	if len(up) != 0 {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(sku.ID)).UpdateSimple(up...)
		if err != nil {
			return nil, s.PushErr(err)
		}

		goodsM := s.GetDao().Good
		_, err = goodsM.WithContext(s.Ctx).Where(goodsM.ID.Eq(sku.GoodsID)).UpdateSimple(goodsM.UpdatedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.GoodsSkuUpdateRes{}, nil
}

// SkuDestroy 删除销售品
func (s Goods) SkuDestroy(req *meet.GoodsSkuDestroyReq) (*meet.GoodsSkuDestroyRes, error) {
	id := utils.StrToLongNumId(req.Id)

	m := s.GetDao().GoodsSku
	sku, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id), m.DeletedAt.Eq(0)).First()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if sku != nil {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(sku.ID)).UpdateSimple(m.DeletedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}

		goodsM := s.GetDao().Good
		_, err = goodsM.WithContext(s.Ctx).Where(goodsM.ID.Eq(sku.GoodsID)).UpdateSimple(goodsM.UpdatedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.GoodsSkuDestroyRes{}, nil
}
