package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"development/stage_five/grpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "无token认证信息")
		}

		var (
			appid  string
			appkey string
		)

		if id, ok := md["appid"]; ok {
			appid = id[0]
		}

		if key, ok := md["appkey"]; ok {
			appkey = key[0]
		}

		if appid != "101010" || appkey != "i am key" {
			return nil, status.Error(codes.Unauthenticated, "无token认证信息")
		}
		fmt.Println("登陆成功！")

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
