package global

import (
	"mxshop-srvs/goods_srv/config"

	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
	EsClient      *elastic.Client
)
