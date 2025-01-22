package admin

import (
	"encoding/json"
	"gorm.io/gen"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/repository/model"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils"
	"jjsdapi/internal/utils/timez"
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
			Or(m.AsTitle.Like(utils.DbAllLike(req.Keyword))).
			Or(m.Code.Like(utils.DbAllLike(req.Keyword))))
	}

	if req.Start != "" {
		conds = append(conds, m.CreatedAt.Gte(timez.TableSearchTime(req.Start)))
	}

	if req.End != "" {
		conds = append(conds, m.CreatedAt.Gte(timez.TableSearchTime(req.End)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PageSize)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.CreatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, s.PushErr(err)
	}

	if total == 0 {
		return &meet.GoodsListRes{}, nil
	}

	var data = make([]*meet.GoodsListOne, 0, len(list))
	for _, v := range list {
		data = append(data, &meet.GoodsListOne{
			Id:          utils.LongNumIdToStr(v.ID),
			Title:       v.Title,
			AsTitle:     v.AsTitle,
			SkuNum:      0,
			Code:        v.Code,
			CreatedDate: timez.TableDateTime(v.CreatedAt),
		})
	}

	return &meet.GoodsListRes{
		List:  data,
		Total: total,
	}, nil
}

// Store 新增商品
func (s Goods) Store(req *meet.GoodsStoreReq) (*meet.GoodsStoreRes, error) {

	skuAttrs, _ := json.Marshal(req.SkuAttrs)
	err := s.GetDao().Good.WithContext(s.Ctx).Create(&model.Good{
		Title:    req.Title,
		AsTitle:  req.AsTitle,
		Code:     req.Code,
		SkuAttrs: string(skuAttrs),
	})
	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.GoodsStoreRes{}, nil
}

// Update 更新商品
func (s Goods) Update(req *meet.GoodsUpdateReq) (*meet.GoodsUpdateRes, error) {
	//TODO 实现业务
	return &meet.GoodsUpdateRes{}, nil
}

// Destroy 删除商品
func (s Goods) Destroy(req *meet.GoodsDestroyReq) (*meet.GoodsDestroyRes, error) {
	//TODO 实现业务
	return &meet.GoodsDestroyRes{}, nil
}
