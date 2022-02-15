package global

import (
	"mxshop-api/user-web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator
)
