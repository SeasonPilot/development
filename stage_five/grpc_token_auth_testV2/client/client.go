package main

import (
	"context"
	"fmt"
	"time"

	"development/stage_five/grpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 拦截器
	inter := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		// metadata
		md := metadata.Pairs("appid", "101010", "appkey", "i am key")
		c := metadata.NewOutgoingContext(ctx, md)

		err := invoker(c, method, req, reply, cc, opts...)
		fmt.Printf("client: %s\n", time.Since(start))
		return err
	}
	opt := grpc.WithUnaryInterceptor(inter)

	// 拨号
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure(), opt)
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
	reply, err := c.SayHello(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reply)
}
