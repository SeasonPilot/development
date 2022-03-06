package global

import (
	"mxshop-api/userop-web/config"
	"mxshop-api/userop-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator

	GoodsClient proto.GoodsClient
)
