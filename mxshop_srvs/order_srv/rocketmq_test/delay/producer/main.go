package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer(primitive.NamesrvAddr{"172.19.30.30:9876"}))
	if err != nil {
		panic(err)
	}

	if err = p.Start(); err != nil {
		panic(err)
	}

	msg := &primitive.Message{
		Topic: "imooc1",
		Body:  []byte("delay message"),
	}
	msg.WithDelayTimeLevel(3)
	resp, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("SendSync err: %s", err)
		return
	} else {
		fmt.Println(resp.String())
	}

	if err = p.Shutdown(); err != nil {
		fmt.Printf("Shutdown err: %s", err)
	}
}
