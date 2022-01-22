package main

import (
	"context"
	"net"

	"development/stage_five/grpc_error_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return nil, status.Errorf(codes.NotFound, "记录未找到:%s", req.Name)
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
