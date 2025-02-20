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
	Id          string                  `json:"id"`           //商品id
	Title       string                  `json:"title"`        //客户名称
	AsTitle     string                  `json:"as_title"`     //商品别名
	SkuNum      int64                   `json:"sku_num"`      //可售商品数量
	UpdatedDate string                  `json:"updated_date"` //更新时间
	GoodsSkus   []*GoodsListOneGoodsSku `json:"goods_skus"`   //可售商品
}

type GoodsListOneGoodsSku struct {
	Id       string `json:"id"`       //商品skuID
	GoodsId  string `json:"goods_id"` //商品ID
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
	Title   string `json:"title"  validate:"trim,required"` //客户名称
	AsTitle string `json:"as_title"  validate:"trim"`       //商品别名
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
	Id      string `json:"id"  validate:"trim,required"`    //商品id
	Title   string `json:"title"  validate:"trim,required"` //客户名称
	AsTitle string `json:"as_title"  validate:"trim"`       //商品别名
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

type GoodsSelectInfoFormat struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data GoodsSelectInfoRes `json:"data"`
}

type GoodsSelectInfoReq struct {
}

type GoodsSelectInfoRes struct {
	Format []*GoodsSelectInfoOne `json:"format"` //商品规格
	Unit   []*GoodsSelectInfoOne `json:"unit"`   //单位
}

type GoodsSelectInfoOne struct {
	Label string `json:"label"` //标签
	Value string `json:"value"` //值
}

// SelectInfo
// @Tags 商品数据
// @Summary 表单要选择的数据
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query GoodsSelectInfoReq true "数据"
// @Success 200 {object} GoodsSelectInfoFormat
// @Router /admin/goods/select-info [GET]
func (Goods) SelectInfo() {

}

type GoodsSkuStoreFormat struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Data GoodsSkuStoreRes `json:"data"`
}

type GoodsSkuStoreReq struct {
	GoodsId  string `json:"goods_id"  validate:"trim,required"` //商品id
	Capacity string `json:"capacity"  validate:"trim"`          //商品重量
	Remark   string `json:"remark"  validate:"trim"`            //商品名称备注
	Format   string `json:"format" `                            //商品规格
	Unit     string `json:"unit" `                              //单位
	Pp       string `json:"pp" `                                //采购价
	Wp       string `json:"wp" `                                //批发价
	Rp       string `json:"rp" `                                //零售价
	Stock    int64  `json:"stock" `                             //库存
	Number   string `json:"number"  validate:"trim"`            //商品编码
}

type GoodsSkuStoreRes struct {
}

// SkuStore
// @Tags 商品数据
// @Summary 创建销售品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsSkuStoreReq true "数据"
// @Success 200 {object} GoodsSkuStoreFormat
// @Router /admin/goods/sku-store [POST]
func (Goods) SkuStore() {

}

type GoodsSkuUpdateFormat struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Data GoodsSkuUpdateRes `json:"data"`
}

type GoodsSkuUpdateReq struct {
	Id       string `json:"id"  validate:"trim,required"` //商品id
	Capacity string `json:"capacity"  validate:"trim"`    //商品重量
	Remark   string `json:"remark"  validate:"trim"`      //商品名称备注
	Format   string `json:"format" `                      //商品规格
	Unit     string `json:"unit" `                        //单位
	Pp       string `json:"pp" `                          //采购价
	Wp       string `json:"wp" `                          //批发价
	Rp       string `json:"rp" `                          //零售价
	Stock    int64  `json:"stock" `                       //库存
	Number   string `json:"number"  validate:"trim"`      //商品编码
}

type GoodsSkuUpdateRes struct {
}

// SkuUpdate
// @Tags 商品数据
// @Summary 更新销售品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsSkuUpdateReq true "数据"
// @Success 200 {object} GoodsSkuUpdateFormat
// @Router /admin/goods/sku-update [POST]
func (Goods) SkuUpdate() {

}

type GoodsSkuDestroyFormat struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data GoodsSkuDestroyRes `json:"data"`
}

type GoodsSkuDestroyReq struct {
	Id string `json:"id"  validate:"trim,required"` //商品id
}

type GoodsSkuDestroyRes struct {
}

// SkuDestroy
// @Tags 商品数据
// @Summary 删除销售品
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data body GoodsSkuDestroyReq true "数据"
// @Success 200 {object} GoodsSkuDestroyFormat
// @Router /admin/goods/sku-destroy [POST]
func (Goods) SkuDestroy() {

}
