package main

import (
	"fmt"
	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	// 创建 Redis 连接
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "172.19.30.30:6379",
	})
	// 创建连接池
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// 从连接池中拿出一个 Redsync 实例
	rs := redsync.New(pool)

	gNum := 2
	mutexname := "my-global-mutex"

	var wg sync.WaitGroup
	for i := 0; i < gNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mutex := rs.NewMutex(mutexname)

			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			// Do your work that requires the lock.
			fmt.Println("获取锁成功")
			time.Sleep(time.Second * 5)

			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
		}()
	}

	wg.Wait()
}
