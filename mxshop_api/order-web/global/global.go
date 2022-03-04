package global

import (
	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator

	GoodsClient proto.GoodsClient

	InvClient proto.InventoryClient

	OrderClient proto.OrderClient
)
