package admin

import (
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
			CustomerId:  utils.LongNumIdToStr(v.CustomerID),
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
	m := s.GetDao().Customer
	list, err := m.WithContext(s.Ctx).Where(m.DeletedAt.Eq(0)).SelectAll()
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

// AddressUpdate 地址更新
func (s Customer) AddressUpdate(req *meet.CustomerAddressUpdateReq) (*meet.CustomerAddressUpdateRes, error) {
	m := s.GetDao().CustomerAddress

	address, err := m.WithContext(s.Ctx).Where(m.ID.Eq(utils.StrToLongNumId(req.Id)), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if address == nil {
		return nil, s.ParamErr("你要找的地址信息不存在")
	}

	var up []field.AssignExpr
	if req.Title != address.Title {
		up = append(up, m.Title.Value(req.Title))
	}

	if req.Address != address.Address {
		up = append(up, m.Address.Value(req.Address))
	}

	if req.Tel != address.Tel {
		up = append(up, m.Tel.Value(req.Tel))
	}

	customerId := utils.StrToLongNumId(req.CustomerId)
	if customerId > 0 && customerId != address.CustomerID {
		up = append(up, m.CustomerID.Value(customerId))
	}

	if len(up) > 0 {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(address.ID)).UpdateSimple(up...)
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.CustomerAddressUpdateRes{}, nil
}

// AddressDestroy 地址删除
func (s Customer) AddressDestroy(req *meet.CustomerAddressDestroyReq) (*meet.CustomerAddressDestroyRes, error) {
	m := s.GetDao().CustomerAddress

	address, err := m.WithContext(s.Ctx).Where(m.ID.Eq(utils.StrToLongNumId(req.Id)), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if address != nil {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(address.ID)).UpdateSimple(m.DeletedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.CustomerAddressDestroyRes{}, nil
}

// Store 新增客户
func (s Customer) Store(req *meet.CustomerStoreReq) (*meet.CustomerStoreRes, error) {
	m := s.GetDao().Customer

	customer, err := m.WithContext(s.Ctx).Where(m.Name.Eq(req.Name), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if customer != nil {
		return nil, s.ParamErr("客户已存在")
	}

	err = m.WithContext(s.Ctx).Create(&model.Customer{
		Name:  req.Name,
		Phone: req.Phone,
	})

	if err != nil {
		return nil, s.PushErr(err)
	}

	return &meet.CustomerStoreRes{}, nil
}

// Update 更新客户
func (s Customer) Update(req *meet.CustomerUpdateReq) (*meet.CustomerUpdateRes, error) {
	id := utils.StrToLongNumId(req.Id)

	m := s.GetDao().Customer
	exist, err := m.WithContext(s.Ctx).Where(m.Name.Eq(req.Name), m.ID.Neq(id), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if exist != nil {
		return nil, s.ParamErr("客户已存在")
	}

	customer, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if customer == nil {
		return nil, s.ParamErr("客户不存在")
	}

	var up []field.AssignExpr
	if req.Name != customer.Name {
		up = append(up, m.Name.Value(req.Name))
	}

	if req.Phone != customer.Phone {
		up = append(up, m.Phone.Value(req.Phone))
	}

	if len(up) > 0 {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(customer.ID)).UpdateSimple(up...)
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.CustomerUpdateRes{}, nil
}

// Destroy 删除客户
func (s Customer) Destroy(req *meet.CustomerDestroyReq) (*meet.CustomerDestroyRes, error) {
	id := utils.StrToLongNumId(req.Id)

	m := s.GetDao().Customer

	customer, err := m.WithContext(s.Ctx).Where(m.ID.Eq(id), m.DeletedAt.Eq(0)).Take()
	if err = dbcheck.DbError(err); err != nil {
		return nil, s.PushErr(err)
	}

	if customer != nil {
		_, err = m.WithContext(s.Ctx).Where(m.ID.Eq(customer.ID)).UpdateSimple(m.DeletedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}

		addressM := s.GetDao().CustomerAddress
		_, err = addressM.WithContext(s.Ctx).Where(addressM.CustomerID.Eq(customer.ID)).UpdateSimple(m.DeletedAt.Value(time.Now().Unix()))
		if err != nil {
			return nil, s.PushErr(err)
		}
	}

	return &meet.CustomerDestroyRes{}, nil
}
