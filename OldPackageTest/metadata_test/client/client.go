package main

import (
	"context"
	"fmt"

	"development/OldPackageTest/grpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	md := metadata.New(map[string]string{
		"name": "season",
	})
	//md1:=metadata.Pairs("hh", "111")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 调用服务
	req := proto.HelloRequest{
		Name: "kk",
	}
	reply, err := c.SayHello(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
