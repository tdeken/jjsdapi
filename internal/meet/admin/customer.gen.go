// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package admin

type CustomerListReq struct {
	Page    int32  `json:"page" query:"page"`                 //当前页
	PerPage int32  `json:"per_page" query:"per_page"`         //每页条数
	Name    string `json:"name" query:"name" validate:"trim"` //客户名称
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

type CustomerAddressListReq struct {
	Page    int32  `json:"page" query:"page"`                   //当前页
	PerPage int32  `json:"per_page" query:"per_page"`           //每页条数
	Title   string `json:"title" query:"title" validate:"trim"` //商店名称
	Tel     string `json:"tel" query:"tel" validate:"trim"`     //联系方式
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
}
