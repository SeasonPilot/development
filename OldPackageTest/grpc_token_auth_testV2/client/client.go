package main

import (
	"context"
	"fmt"

	"development/OldPackageTest/grpc_token_auth_testV2/proto"

	"google.golang.org/grpc"
)

type customCredential struct{}

func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	fmt.Println(ctx, uri)
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

func (c *customCredential) RequireTransportSecurity() bool {
	return false
}

func main() {
	// 拨号
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&customCredential{}))
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
