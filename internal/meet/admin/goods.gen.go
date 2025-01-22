// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package admin

type GoodsListReq struct {
	Page     int32  `json:"page" query:"page"`                                                     //当前页
	PageSize int32  `json:"page_size" query:"page_size"`                                           //每页条数
	Keyword  string `json:"keyword" query:"keyword" validate:"trim"`                               //关键字
	Start    string `json:"start" query:"start" validate:"omitempty,datetime=2006-01-02 15:04:05"` //关键字
	End      string `json:"end" query:"end" validate:"omitempty,datetime=2006-01-02 15:04:05"`     //关键字
}

type GoodsListRes struct {
	List  []*GoodsListOne `json:"list"`  //列表数据
	Total int64           `json:"total"` //数据总条数
}

type GoodsListOne struct {
	Id          string `json:"id"`           //商品id
	Title       string `json:"title"`        //客户名称
	AsTitle     string `json:"as_title"`     //商品别名
	SkuNum      int64  `json:"sku_num"`      //可售商品数量
	Code        string `json:"code"`         //商品编号
	CreatedDate string `json:"created_date"` //创建时间
}

type GoodsStoreReq struct {
	Title    string              `json:"title" form:"title" validate:"trim,required"`         //客户名称
	AsTitle  string              `json:"as_title" form:"as_title" validate:"trim"`            //商品别名
	Code     string              `json:"code" form:"code" validate:"trim"`                    //商品别名
	SkuAttrs [][]*GoodsStoreAttr `json:"sku_attrs" form:"sku_attrs" validate:"required,give"` //商品sku属性
}

type GoodsStoreAttr struct {
	Mark     string `json:"mark" form:"mark" validate:"required"`              //属性
	ShowType int32  `json:"show_type" form:"show_type" validate:"oneof=1 2 3"` //展示方式(1-不展示，2带括号，3-不带括号)
}

type GoodsStoreRes struct {
}

type GoodsUpdateReq struct {
	Id       string              `json:"id" form:"id" validate:"trim,required"`               //商品id
	Title    string              `json:"title" form:"title" validate:"trim,required"`         //客户名称
	AsTitle  string              `json:"as_title" form:"as_title" validate:"trim"`            //商品别名
	Code     string              `json:"code" form:"code" validate:"trim"`                    //商品别名
	SkuAttrs []*GoodsUpdateAttrs `json:"sku_attrs" form:"sku_attrs" validate:"required,give"` //商品sku属性
}

type GoodsUpdateAttrs struct {
	Attrs []*GoodsUpdateAttrsAttr `json:"attrs" form:"attrs" validate:"required,give"` //属性集合
}

type GoodsUpdateAttrsAttr struct {
	Mark     string `json:"mark" form:"mark" validate:"required"`              //属性
	ShowType int32  `json:"show_type" form:"show_type" validate:"oneof=1 2 3"` //展示方式(1-不展示，2带括号，3-不带括号)
}

type GoodsUpdateRes struct {
}

type GoodsDestroyReq struct {
	Id string `json:"id" form:"id" validate:"trim,required"` //商品id
}

type GoodsDestroyRes struct {
}
