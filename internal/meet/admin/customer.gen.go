// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package admin

type CustomerAddressListReq struct {
	Page    int32 `json:"page" query:"page"`         //当前页
	PerPage int32 `json:"per_page" query:"per_page"` //每页条数
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
