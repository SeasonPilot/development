package global

import (
	"mxshop-srvs/order_srv/config"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig config.ServerConfig
	RedClient     *redis.Client
	RedSync       *redsync.Redsync
)
