package admin

import (
	"gorm.io/gen"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils"
	"time"
)

type Customer struct {
	service.Service
}

// AddressList 配送地址列表
func (s Customer) AddressList(req *meet.CustomerAddressListReq) (*meet.CustomerAddressListRes, error) {
	m := s.GetDao().CustomerAddress

	var conds []gen.Condition

	offset, limit := utils.PpDbLo(req.Page, req.PerPage)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.ID.Desc()).FindByPage(offset, limit)
	if err != nil || total == 0 {
		return nil, s.PushErr(err)
	}

	var res = &meet.CustomerAddressListRes{
		List:  make([]*meet.CustomerAddressListOne, 0, len(list)),
		Total: total,
	}

	for _, v := range list {
		res.List = append(res.List, &meet.CustomerAddressListOne{
			Id:          v.ID,
			Title:       v.Title,
			Address:     v.Address,
			Tel:         v.Tel,
			CreatedDate: time.Unix(v.CreatedAt, 0).Format(time.DateTime),
		})
	}

	return res, nil
}
