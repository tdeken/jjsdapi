// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package admin

type CustomerListReq struct {
	Page     int32  `json:"page" query:"page"`                 //当前页
	PageSize int32  `json:"page_size" query:"page_size"`       //每页条数
	Name     string `json:"name" query:"name" validate:"trim"` //客户名称
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

type CustomerSelectReq struct {
}

type CustomerSelectRes struct {
	List []*CustomerSelectOne `json:"list"` //列表数据
}

type CustomerSelectOne struct {
	Id   string `json:"id"`   //地址id
	Name string `json:"name"` //客户名称
}

type CustomerStoreReq struct {
	Name  string `json:"name" form:"name" validate:"trim,required"` //客户名称
	Phone string `json:"phone" form:"phone" validate:"trim"`        //客户手机号
}

type CustomerStoreRes struct {
}

type CustomerUpdateReq struct {
	Id    string `json:"id" form:"id" validate:"required"`          //客户id
	Name  string `json:"name" form:"name" validate:"trim,required"` //客户名称
	Phone string `json:"phone" form:"phone" validate:"trim"`        //客户手机号
}

type CustomerUpdateRes struct {
}

type CustomerDestroyReq struct {
	Id string `json:"id" form:"id" validate:"required"` //客户id
}

type CustomerDestroyRes struct {
}

type CustomerAddressListReq struct {
	Page     int32  `json:"page" query:"page"`                   //当前页
	PageSize int32  `json:"page_size" query:"page_size"`         //每页条数
	Title    string `json:"title" query:"title" validate:"trim"` //商店名称
	Tel      string `json:"tel" query:"tel" validate:"trim"`     //联系方式
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

type CustomerAddressCreateReq struct {
	Title      string `json:"title" form:"title" validate:"trim,required"`     //商店名称
	Address    string `json:"address" form:"address" validate:"trim,required"` //商店地址
	Tel        string `json:"tel" form:"tel" validate:"trim"`                  //联系方式
	CustomerId string `json:"customer_id" form:"customer_id" validate:"trim"`  //客户id
}

type CustomerAddressCreateRes struct {
}

type CustomerAddressUpdateReq struct {
	Id         string `json:"id" form:"id" validate:"required"`                //商店地址id
	Title      string `json:"title" form:"title" validate:"trim,required"`     //商店名称
	Address    string `json:"address" form:"address" validate:"trim,required"` //商店地址
	Tel        string `json:"tel" form:"tel" validate:"trim"`                  //联系方式
	CustomerId string `json:"customer_id" form:"customer_id" validate:"trim"`  //客户id
}

type CustomerAddressUpdateRes struct {
}

type CustomerAddressDestroyReq struct {
	Id string `json:"id" form:"id" validate:"required"` //商店地址id
}

type CustomerAddressDestroyRes struct {
}
