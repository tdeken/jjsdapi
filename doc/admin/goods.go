package admin

type Goods struct {
}

type GoodsListFormat struct {
	Code int32        `json:"code"`
	Msg  string       `json:"msg"`
	Data GoodsListRes `json:"data"`
}

type GoodsListReq struct {
	Page     int32  `json:"page" `                                                    //当前页
	PageSize int32  `json:"page_size" `                                               //每页条数
	Keyword  string `json:"keyword"  validate:"trim"`                                 //关键字
	Start    string `json:"start"  validate:"omitempty,datetime=2006-01-02 15:04:05"` //关键字
	End      string `json:"end"  validate:"omitempty,datetime=2006-01-02 15:04:05"`   //关键字
}

type GoodsListRes struct {
	List  []*GoodsListOne `json:"list"`  //列表数据
	Total int64           `json:"total"` //数据总条数
}

type GoodsListOne struct {
	Id          string                `json:"id"`           //商品id
	Title       string                `json:"title"`        //客户名称
	AsTitle     string                `json:"as_title"`     //商品别名
	SkuNum      int64                 `json:"sku_num"`      //可售商品数量
	Code        string                `json:"code"`         //商品编号
	CreatedDate string                `json:"created_date"` //创建时间
	SkuAttrs    [][]*GoodsListOneAttr `json:"sku_attrs"`    //商品sku属性
}

type GoodsListOneAttr struct {
	Mark     string `json:"mark"`      //属性
	ShowType int32  `json:"show_type"` //展示方式(1-不展示，2带括号，3-不带括号)
}

// List
// @Tags 商品数据
// @Summary 客户列表
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query GoodsListReq true "数据"
// @Success 200 {object} GoodsListFormat
// @Router /admin/goods/list [GET]
func (Goods) List() {

}

type GoodsStoreFormat struct {
	Code int32         `json:"code"`
	Msg  string        `json:"msg"`
	Data GoodsStoreRes `json:"data"`
}

type GoodsStoreReq struct {
	Title    string              `json:"title"  validate:"trim,required"`     //客户名称
	AsTitle  string              `json:"as_title"  validate:"trim"`           //商品别名
	Code     string              `json:"code"  validate:"trim"`               //商品别名
	SkuAttrs [][]*GoodsStoreAttr `json:"sku_attrs"  validate:"required,give"` //商品sku属性
}

type GoodsStoreAttr struct {
	Mark     string `json:"mark"  validate:"required"`         //属性
	ShowType int32  `json:"show_type"  validate:"oneof=1 2 3"` //展示方式(1-不展示，2带括号，3-不带括号)
}

type GoodsStoreRes struct {
}

// Store
// @Tags 商品数据
// @Summary 新增商品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsStoreReq true "数据"
// @Success 200 {object} GoodsStoreFormat
// @Router /admin/goods/store [POST]
func (Goods) Store() {

}

type GoodsUpdateFormat struct {
	Code int32          `json:"code"`
	Msg  string         `json:"msg"`
	Data GoodsUpdateRes `json:"data"`
}

type GoodsUpdateReq struct {
	Id       string              `json:"id"  validate:"trim,required"`        //商品id
	Title    string              `json:"title"  validate:"trim,required"`     //客户名称
	AsTitle  string              `json:"as_title"  validate:"trim"`           //商品别名
	Code     string              `json:"code"  validate:"trim"`               //商品别名
	SkuAttrs []*GoodsUpdateAttrs `json:"sku_attrs"  validate:"required,give"` //商品sku属性
}

type GoodsUpdateAttrs struct {
	Attrs []*GoodsUpdateAttrsAttr `json:"attrs"  validate:"required,give"` //属性集合
}

type GoodsUpdateAttrsAttr struct {
	Mark     string `json:"mark"  validate:"required"`         //属性
	ShowType int32  `json:"show_type"  validate:"oneof=1 2 3"` //展示方式(1-不展示，2带括号，3-不带括号)
}

type GoodsUpdateRes struct {
}

// Update
// @Tags 商品数据
// @Summary 更新商品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsUpdateReq true "数据"
// @Success 200 {object} GoodsUpdateFormat
// @Router /admin/goods/update [POST]
func (Goods) Update() {

}

type GoodsDestroyFormat struct {
	Code int32           `json:"code"`
	Msg  string          `json:"msg"`
	Data GoodsDestroyRes `json:"data"`
}

type GoodsDestroyReq struct {
	Id string `json:"id"  validate:"trim,required"` //商品id
}

type GoodsDestroyRes struct {
}

// Destroy
// @Tags 商品数据
// @Summary 删除商品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsDestroyReq true "数据"
// @Success 200 {object} GoodsDestroyFormat
// @Router /admin/goods/destroy [POST]
func (Goods) Destroy() {

}
