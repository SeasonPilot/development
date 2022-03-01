package initialization

import (
	"context"
	"fmt"

	"mxshop-srvs/inventory_srv/global"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func InitRedisClient() {
	global.RedClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServiceConfig.RedisInfo.Host, global.ServiceConfig.RedisInfo.Port),
	})

	s := global.RedClient.Ping(context.Background())
	if s.Err() != nil {
		panic(s.Err())
	}
}

func InitRedSync() {
	// 创建连接池
	pool := goredis.NewPool(global.RedClient) // or, pool := redigo.NewPool(...)

	// 从连接池中拿出一个 Redsync 实例
	global.RedSync = redsync.New(pool)
}
