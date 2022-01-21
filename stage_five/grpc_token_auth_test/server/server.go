package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"development/stage_five/grpc_test/proto"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
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

	// 拦截器
	inter := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		fmt.Printf("server: %s\n", time.Since(start))
		return resp, err
	}
	opt := grpc.UnaryInterceptor(inter)

	// 注册服务
	g := grpc.NewServer(opt)
	proto.RegisterGreeterServer(g, &Server{})
	// 启动服务
	err = g.Serve(listener)
	if err != nil {
		panic(err)
	}

}
