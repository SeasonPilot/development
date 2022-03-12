package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct{}

func (*OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	<-time.After(time.Second * 3)
	fmt.Println("执行本地逻辑失败")
	return primitive.UnknowState
}

func (*OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	<-time.After(time.Second * 15)
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(&OrderListener{},
		producer.WithNameServer(primitive.NamesrvAddr{"172.19.30.30:9876"}))
	if err != nil {
		panic(err)
	}

	if err = p.Start(); err != nil {
		panic(err)
	}

	resp, err := p.SendMessageInTransaction(context.Background(), &primitive.Message{
		Topic: "trans topic",
		Body:  []byte("this is transaction message2"),
	})
	if err != nil {
		fmt.Printf("SendMessageInTransaction err: %s", err)
		return
	} else {
		fmt.Println(resp.String())
	}

	<-time.After(time.Hour)
	if err = p.Shutdown(); err != nil {
		fmt.Printf("Shutdown err: %s", err)
	}
}
