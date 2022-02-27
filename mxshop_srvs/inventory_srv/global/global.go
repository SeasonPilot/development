package global

import (
	"mxshop-srvs/inventory_srv/config"

	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
)
