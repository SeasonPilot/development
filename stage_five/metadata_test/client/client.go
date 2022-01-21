package main

import (
	"context"
	"fmt"

	"development/stage_five/grpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 拨号
	conn, err := grpc.Dial("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// new client 对象
	c := proto.NewGreeterClient(conn)

	// 调用服务
	req := proto.HelloRequest{
		Name: "kk",
	}
	reply, err := c.SayHello(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
