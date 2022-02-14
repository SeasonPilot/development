package main

import (
	"context"
	"fmt"
	"net"

	"development/OldPackageTest/grpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	// metadata.FromIncomingContext 是写在业务逻辑中，不是 main函数中
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md["name"])
	}

	return &proto.HelloReply{
		Massage: "Hello," + req.Name,
	}, nil
}

func main() {
	// 监听
	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}

	// 注册服务
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	// 启动服务
	err = g.Serve(listener)
	if err != nil {
		panic(err)
	}

}
