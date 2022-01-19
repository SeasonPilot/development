package main

import (
	"development/stage_five/new_helloworld/handler"
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}

	var reply string
	err = client.Call(handler.HelloServiceName+".Hello", "season", &reply)
	if err != nil {
		panic("调用失败")
	}

	fmt.Println(reply)
}
