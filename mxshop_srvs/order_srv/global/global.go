package global

import (
	"mxshop-srvs/order_srv/config"
	"mxshop-srvs/order_srv/proto"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	ServiceConfig   config.ServerConfig
	RedClient       *redis.Client
	RedSync         *redsync.Redsync
	GoodsClient     proto.GoodsClient
	InventoryClient proto.InventoryClient
)
