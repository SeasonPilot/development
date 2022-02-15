package global

import (
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator

	UserClient proto.UserClient
)
