package global

import (
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator

	GoodsClient proto.GoodsClient
)
