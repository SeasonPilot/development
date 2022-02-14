package main

import (
	"context"
	"net"
	"time"

	"development/OldPackageTest/grpc_overtime_test/proto"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	time.Sleep(time.Second * 5)
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
