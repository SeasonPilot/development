package global

import (
	"mxshop-srvs/userop_srv/config"

	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
)
