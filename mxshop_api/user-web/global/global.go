package global

import (
	"development/mxshop_api/user-web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	SrvConfig = &config.ServerConfig{}

	Translator ut.Translator
)
