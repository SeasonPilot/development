package global

import (
	"mxshop-srvs/user_srv/config"

	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
)
