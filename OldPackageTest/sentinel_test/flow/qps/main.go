package main

import (
	"fmt"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
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
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}

	for i := 0; i < 12; i++ {
		// 3. Entry 方法用于埋点
		e, b := sentinel.Entry("some-test")
		if b != nil {
			fmt.Printf("限流了 %v\n", b.Error())
		} else {
			fmt.Println("检查通过")
			// 务必保证业务逻辑结束后 Exit
			e.Exit()
		}
	}
}
