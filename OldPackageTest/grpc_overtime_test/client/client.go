package main

import (
	"context"
	"fmt"
	"time"

	"development/OldPackageTest/grpc_overtime_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	// 拨号
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	//grpc.Dial(":9950", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	_, err = c.SayHello(ctx, &req)
	if err != nil {
		sts, ok := status.FromError(err)
		if !ok {
			// Error was not a status error
			panic("解析error失败")
		}
		fmt.Println(sts.Message())
		fmt.Println(sts.Code())
	}
}
