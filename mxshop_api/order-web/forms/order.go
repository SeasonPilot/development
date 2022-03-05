package forms

type ShopCartItemForm struct {
	GoodsId int32 `json:"goods_id" form:"goods_id" binding:"required"`
	Num     int32 `json:"num" form:"num" binding:"required"`
}
