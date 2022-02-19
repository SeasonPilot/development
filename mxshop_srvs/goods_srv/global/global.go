package global

import (
	"mxshop-srvs/goods_srv/config"

	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
)
