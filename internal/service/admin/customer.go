package admin

import (
	"gorm.io/gen"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils"
	"jjsdapi/internal/utils/timez"
)

type Customer struct {
	service.Service
}

// AddressList 配送地址列表
func (s Customer) AddressList(req *meet.CustomerAddressListReq) (*meet.CustomerAddressListRes, error) {
	m := s.GetDao().CustomerAddress

	var conds []gen.Condition
	if req.Title != "" {
		conds = append(conds, m.Title.Like(utils.DbAllLike(req.Title)))
	}

	if req.Tel != "" {
		conds = append(conds, m.Tel.Like(utils.DbAllLike(req.Tel)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PerPage)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.CreatedAt.Desc()).FindByPage(offset, limit)
	if err != nil || total == 0 {
		return nil, s.PushErr(err)
	}

	var res = &meet.CustomerAddressListRes{
		List:  make([]*meet.CustomerAddressListOne, 0, len(list)),
		Total: total,
	}

	for _, v := range list {
		res.List = append(res.List, &meet.CustomerAddressListOne{
			Id:          utils.LongNumIdToStr(v.ID),
			Title:       v.Title,
			Address:     v.Address,
			Tel:         v.Tel,
			CreatedDate: timez.TableDateTime(v.CreatedAt),
		})
	}

	return res, nil
}

// List 客户列表
func (s Customer) List(req *meet.CustomerListReq) (*meet.CustomerListRes, error) {
	m := s.GetDao().Customer

	var conds []gen.Condition
	if req.Name != "" {
		conds = append(conds, m.Name.Like(utils.DbAllLike(req.Name)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PerPage)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.CreatedAt.Desc()).FindByPage(offset, limit)
	if err != nil || total == 0 {
		return nil, s.PushErr(err)
	}

	var res = &meet.CustomerListRes{
		List:  make([]*meet.CustomerListOne, 0, len(list)),
		Total: total,
	}

	for _, v := range list {
		res.List = append(res.List, &meet.CustomerListOne{
			Id:          utils.LongNumIdToStr(v.ID),
			Name:        v.Name,
			Phone:       v.Phone,
			CreatedDate: timez.TableDateTime(v.CreatedAt),
		})
	}

	return res, nil
}
