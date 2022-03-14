package main

import (
	"context"
	"fmt"
	"time"

	"development/OldPackageTest/grpc_test/proto"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	// 拦截器
	inter := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("client: %s\n", time.Since(start))
		return err
	}
	opt := grpc.WithUnaryInterceptor(inter)
	retryOpts := grpc.WithUnaryInterceptor(
		grpc_retry.UnaryClientInterceptor(
			// 重试次数
			grpc_retry.WithMax(3),
			// 超时时间
			grpc_retry.WithPerRetryTimeout(2*time.Second),
			// 设置应该重试哪些状态码
			grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
		),
	)
	// 拨号
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure(), opt, retryOpts)
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
