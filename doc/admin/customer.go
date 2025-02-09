package admin

type Customer struct {
}

type CustomerListFormat struct {
	Code int32           `json:"code"`
	Msg  string          `json:"msg"`
	Data CustomerListRes `json:"data"`
}

type CustomerListReq struct {
	Page     int32  `json:"page" `                 //当前页
	PageSize int32  `json:"page_size" `            //每页条数
	Name     string `json:"name"  validate:"trim"` //客户名称
}

type CustomerListRes struct {
	List  []*CustomerListOne `json:"list"`  //列表数据
	Total int64              `json:"total"` //数据总条数
}

type CustomerListOne struct {
	Id          string `json:"id"`           //地址id
	Name        string `json:"name"`         //客户名称
	Phone       string `json:"phone"`        //客户手机号
	CreatedDate string `json:"created_date"` //创建时间
}

// List
// @Tags 客户数据
// @Summary 客户列表
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query CustomerListReq true "数据"
// @Success 200 {object} CustomerListFormat
// @Router /admin/customer/list [GET]
func (Customer) List() {

}

type CustomerSelectFormat struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Data CustomerSelectRes `json:"data"`
}

type CustomerSelectReq struct {
}

type CustomerSelectRes struct {
	List []*CustomerSelectOne `json:"list"` //列表数据
}

type CustomerSelectOne struct {
	Id   string `json:"id"`   //地址id
	Name string `json:"name"` //客户名称
}

// Select
// @Tags 客户数据
// @Summary 客户列表选择
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query CustomerSelectReq true "数据"
// @Success 200 {object} CustomerSelectFormat
// @Router /admin/customer/select [GET]
func (Customer) Select() {

}

type CustomerStoreFormat struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Data CustomerStoreRes `json:"data"`
}

type CustomerStoreReq struct {
	Name  string `json:"name"  validate:"trim,required"` //客户名称
	Phone string `json:"phone"  validate:"trim"`         //客户手机号
}

type CustomerStoreRes struct {
}

// Store
// @Tags 客户数据
// @Summary 新增客户
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerStoreReq true "数据"
// @Success 200 {object} CustomerStoreFormat
// @Router /admin/customer/store [POST]
func (Customer) Store() {

}

type CustomerUpdateFormat struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Data CustomerUpdateRes `json:"data"`
}

type CustomerUpdateReq struct {
	Id    string `json:"id"  validate:"required"`        //客户id
	Name  string `json:"name"  validate:"trim,required"` //客户名称
	Phone string `json:"phone"  validate:"trim"`         //客户手机号
}

type CustomerUpdateRes struct {
}

// Update
// @Tags 客户数据
// @Summary 更新客户
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerUpdateReq true "数据"
// @Success 200 {object} CustomerUpdateFormat
// @Router /admin/customer/update [POST]
func (Customer) Update() {

}

type CustomerDestroyFormat struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data CustomerDestroyRes `json:"data"`
}

type CustomerDestroyReq struct {
	Id string `json:"id"  validate:"required"` //客户id
}

type CustomerDestroyRes struct {
}

// Destroy
// @Tags 客户数据
// @Summary 删除客户
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerDestroyReq true "数据"
// @Success 200 {object} CustomerDestroyFormat
// @Router /admin/customer/destroy [POST]
func (Customer) Destroy() {

}

type CustomerAddressListFormat struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data CustomerAddressListRes `json:"data"`
}

type CustomerAddressListReq struct {
	Page     int32  `json:"page" `                  //当前页
	PageSize int32  `json:"page_size" `             //每页条数
	Title    string `json:"title"  validate:"trim"` //商店名称
	Tel      string `json:"tel"  validate:"trim"`   //联系方式
}

type CustomerAddressListRes struct {
	List  []*CustomerAddressListOne `json:"list"`  //列表数据
	Total int64                     `json:"total"` //数据总条数
}

type CustomerAddressListOne struct {
	Id          string `json:"id"`           //地址id
	Title       string `json:"title"`        //商店名称
	Address     string `json:"address"`      //商店地址
	Tel         string `json:"tel"`          //联系电话
	CreatedDate string `json:"created_date"` //创建时间
	CustomerId  string `json:"customer_id"`  //客户id
}

// AddressList
// @Tags 客户数据
// @Summary 配送地址列表
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query CustomerAddressListReq true "数据"
// @Success 200 {object} CustomerAddressListFormat
// @Router /admin/customer/address-list [GET]
func (Customer) AddressList() {

}

type CustomerAddressCreateFormat struct {
	Code int32                    `json:"code"`
	Msg  string                   `json:"msg"`
	Data CustomerAddressCreateRes `json:"data"`
}

type CustomerAddressCreateReq struct {
	Title      string `json:"title"  validate:"trim,required"`   //商店名称
	Address    string `json:"address"  validate:"trim,required"` //商店地址
	Tel        string `json:"tel"  validate:"trim"`              //联系方式
	CustomerId string `json:"customer_id"  validate:"trim"`      //客户id
}

type CustomerAddressCreateRes struct {
}

// AddressCreate
// @Tags 客户数据
// @Summary 地址创建
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerAddressCreateReq true "数据"
// @Success 200 {object} CustomerAddressCreateFormat
// @Router /admin/customer/address-create [POST]
func (Customer) AddressCreate() {

}

type CustomerAddressUpdateFormat struct {
	Code int32                    `json:"code"`
	Msg  string                   `json:"msg"`
	Data CustomerAddressUpdateRes `json:"data"`
}

type CustomerAddressUpdateReq struct {
	Id         string `json:"id"  validate:"required"`           //商店地址id
	Title      string `json:"title"  validate:"trim,required"`   //商店名称
	Address    string `json:"address"  validate:"trim,required"` //商店地址
	Tel        string `json:"tel"  validate:"trim"`              //联系方式
	CustomerId string `json:"customer_id"  validate:"trim"`      //客户id
}

type CustomerAddressUpdateRes struct {
}

// AddressUpdate
// @Tags 客户数据
// @Summary 地址更新
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerAddressUpdateReq true "数据"
// @Success 200 {object} CustomerAddressUpdateFormat
// @Router /admin/customer/address-update [POST]
func (Customer) AddressUpdate() {

}

type CustomerAddressDestroyFormat struct {
	Code int32                     `json:"code"`
	Msg  string                    `json:"msg"`
	Data CustomerAddressDestroyRes `json:"data"`
}

type CustomerAddressDestroyReq struct {
	Id string `json:"id"  validate:"required"` //商店地址id
}

type CustomerAddressDestroyRes struct {
}

// AddressDestroy
// @Tags 客户数据
// @Summary 地址删除
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body CustomerAddressDestroyReq true "数据"
// @Success 200 {object} CustomerAddressDestroyFormat
// @Router /admin/customer/address-destroy [POST]
func (Customer) AddressDestroy() {

}
