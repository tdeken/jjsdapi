package admin

import (
	"gorm.io/gen"
	meet "jjsdapi/internal/meet/admin"
	"jjsdapi/internal/repository/model"
	"jjsdapi/internal/service"
	"jjsdapi/internal/utils"
	"jjsdapi/internal/utils/dbcheck"
	"jjsdapi/internal/utils/timez"
)

type Customer struct {
	service.Service
}

// AddressList 配送地址列表
func (s Customer) AddressList(req *meet.CustomerAddressListReq) (*meet.CustomerAddressListRes, error) {
	m := s.GetDao().CustomerAddress

	var conds = []gen.Condition{
		m.DeletedAt.Eq(0),
	}

	if req.Title != "" {
		conds = append(conds, m.Title.Like(utils.DbAllLike(req.Title)))
	}

	if req.Tel != "" {
		conds = append(conds, m.Tel.Like(utils.DbAllLike(req.Tel)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PageSize)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.CreatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
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

	var conds = []gen.Condition{
		m.DeletedAt.Eq(0),
	}

	if req.Name != "" {
		conds = append(conds, m.Name.Like(utils.DbAllLike(req.Name)))
	}

	offset, limit := utils.PpDbLo(req.Page, req.PageSize)
	list, total, err := m.WithContext(s.Ctx).Where(conds...).Order(m.CreatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
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

// Select 客户列表选择
func (s Customer) Select(req *meet.CustomerSelectReq) (*meet.CustomerSelectRes, error) {
	list, err := s.GetDao().Customer.WithContext(s.Ctx).SelectAll()
	if err != nil {
		return nil, s.PushErr(err)
	}

	var res = &meet.CustomerSelectRes{
		List: make([]*meet.CustomerSelectOne, 0, len(list)),
	}

	for _, v := range list {
		res.List = append(res.List, &meet.CustomerSelectOne{
			Id:   utils.LongNumIdToStr(v.ID),
			Name: v.Name,
		})
	}

	return res, nil
}

// AddressCreate 地址创建
func (s Customer) AddressCreate(req *meet.CustomerAddressCreateReq) (*meet.CustomerAddressCreateRes, error) {
	m := s.GetDao().CustomerAddress

	customerId := utils.StrToLongNumId(req.CustomerId)
	if customerId > 0 {

		list, err := m.WithContext(s.Ctx).Where(m.CustomerID.Eq(customerId), m.DeletedAt.Eq(0)).Find()
		if err != nil {
			return nil, s.PushErr(err)
		}

		for _, v := range list {
			if v.Address == req.Address && v.Title == req.Title {
				return nil, s.ParamErr("配送地址已存在")
			}
		}
	} else {
		customerM := s.GetDao().Customer
		customer, err := customerM.WithContext(s.Ctx).Where(customerM.Name.Eq(req.Title), customerM.DeletedAt.Eq(0)).Take()
		if err = dbcheck.DbError(err); err != nil {
			return nil, s.PushErr(err)
		}

		if customer == nil {
			customer = &model.Customer{
				Name:  req.Title,
				Phone: req.Tel,
			}
			err = customerM.WithContext(s.Ctx).Create(customer)
			if err != nil {
				return nil, s.PushErr(err)
			}
		}
		customerId = customer.ID
	}

	err := m.WithContext(s.Ctx).Create(&model.CustomerAddress{
		CustomerID: customerId,
		Title:      req.Title,
		Address:    req.Address,
		Tel:        req.Tel,
	})

	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.CustomerAddressCreateRes{}, nil
}
