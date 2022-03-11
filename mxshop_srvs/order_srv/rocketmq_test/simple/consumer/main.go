package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(primitive.NamesrvAddr{"172.19.30.30:9876"}),
		// 通过 GroupName 可以达到负载均衡的效果
		consumer.WithGroupName("mxshop"),
	)

	if err := c.Subscribe("imooc1", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range ext {
				fmt.Printf("获得的消息: %v", msg)
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		fmt.Println("获得消息失败")
		return
	}

	_ = c.Start()

	// 不能让主goroutine退出
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}
