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
