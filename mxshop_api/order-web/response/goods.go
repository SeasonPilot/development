package response

import "mxshop-api/order-web/proto"

type GoodsForm struct {
	Id              int32     `json:"id,omitempty"`
	CategoryId      int32     `json:"category_id,omitempty"`
	Name            string    `json:"name,omitempty"`
	GoodsSn         string    `json:"goods_sn,omitempty"`
	ClickNum        int32     `json:"click_num,omitempty"`
	SoldNum         int32     `json:"sold_num,omitempty"`
	FavNum          int32     `json:"fav_num,omitempty"`
	MarketPrice     float32   `json:"market_price,omitempty"`
	ShopPrice       float32   `json:"shop_price,omitempty"`
	GoodsBrief      string    `json:"goods_brief,omitempty"`
	GoodsDesc       string    `json:"goods_desc,omitempty"`
	ShipFree        bool      `json:"ship_free,omitempty"`
	Images          []string  `json:"images,omitempty"`
	DescImages      []string  `json:"desc_images,omitempty"`
	GoodsFrontImage string    `json:"goods_front_image,omitempty"`
	IsNew           bool      `json:"is_new,omitempty"`
	IsHot           bool      `json:"is_hot,omitempty"`
	OnSale          bool      `json:"on_sale,omitempty"`
	AddTime         int64     `json:"add_time,omitempty"`
	Category        *Category `json:"category,omitempty"`
	Brand           *Brand    `json:"brand"`
}
type Category struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Brand struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Logo string `json:"logo,omitempty"`
}

func RespToModels(rsp []*proto.GoodsInfoResponse) []GoodsForm {
	var gf []GoodsForm
	for _, r := range rsp {
		gf = append(gf, GoodsForm{
			Id:              r.Id,
			CategoryId:      r.CategoryId,
			Name:            r.Name,
			GoodsSn:         r.GoodsSn,
			ClickNum:        r.ClickNum,
			SoldNum:         r.SoldNum,
			FavNum:          r.FavNum,
			MarketPrice:     r.MarketPrice,
			ShopPrice:       r.ShopPrice,
			GoodsBrief:      r.GoodsBrief,
			GoodsDesc:       r.GoodsDesc,
			ShipFree:        r.ShipFree,
			Images:          r.Images,
			DescImages:      r.DescImages,
			GoodsFrontImage: r.GoodsFrontImage,
			IsNew:           r.IsNew,
			IsHot:           r.IsHot,
			OnSale:          r.OnSale,
			AddTime:         r.AddTime,
			Category: &Category{
				Id:   r.Category.Id,
				Name: r.Category.Name,
			},
			Brand: &Brand{
				Id:   r.Brand.Id,
				Name: r.Brand.Name,
				Logo: r.Brand.Logo,
			},
		})
	}
	return gf
}
