package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	// 1.先初始化 sentinel
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatal(err)
	}

	// 2.配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              1000,
			WarmUpPeriodSec:        30,
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}

	var globalTotal int
	var passTotal int
	var blockTotal int

	// 阻塞主协程退出
	ch := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func() {
			for {
				globalTotal++
				// 3. Entry 方法用于埋点
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					blockTotal++
					<-time.After(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					<-time.After(time.Duration(rand.Uint64()%10) * time.Millisecond)
					// 务必保证业务逻辑结束后 Exit
					e.Exit()
				}
			}
		}()
	}

	go func() {
		var oldTotal int //过去1s总共有多少个
		var oldPass int  //过去1s总共pass多少个
		var oldBlock int //过去1s总共block多少个
		for {
			oneSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal

			<-time.After(time.Second)
			fmt.Printf("total:%d, pass:%d, block:%d\n", oneSecondTotal, oneSecondPass, oneSecondBlock)
		}
	}()

	<-ch
}
