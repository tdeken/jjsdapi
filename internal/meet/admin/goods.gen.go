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
	Id          string                  `json:"id"`           //商品id
	Title       string                  `json:"title"`        //客户名称
	AsTitle     string                  `json:"as_title"`     //商品别名
	SkuNum      int64                   `json:"sku_num"`      //可售商品数量
	UpdatedDate string                  `json:"updated_date"` //更新时间
	GoodsSkus   []*GoodsListOneGoodsSku `json:"goods_skus"`   //可售商品
}

type GoodsListOneGoodsSku struct {
	Id       string `json:"id"`       //商品skuID
	Name     string `json:"name"`     //销售商品名称
	Capacity string `json:"capacity"` //商品容量
	Remark   string `json:"remark"`   //商品备注
	Format   string `json:"format"`   //规格
	Unit     string `json:"unit"`     //单位
	Pp       string `json:"pp"`       //采购价
	Wp       string `json:"wp"`       //批发价
	Rp       string `json:"rp"`       //零售价
	Stock    int64  `json:"stock"`    //库存
	Number   string `json:"number"`   //商品编号
}

type GoodsStoreReq struct {
	Title   string `json:"title" form:"title" validate:"trim,required"` //客户名称
	AsTitle string `json:"as_title" form:"as_title" validate:"trim"`    //商品别名
}

type GoodsStoreRes struct {
}

type GoodsUpdateReq struct {
	Id      string `json:"id" form:"id" validate:"trim,required"`       //商品id
	Title   string `json:"title" form:"title" validate:"trim,required"` //客户名称
	AsTitle string `json:"as_title" form:"as_title" validate:"trim"`    //商品别名
}

type GoodsUpdateRes struct {
}

type GoodsDestroyReq struct {
	Id string `json:"id" form:"id" validate:"trim,required"` //商品id
}

type GoodsDestroyRes struct {
}
