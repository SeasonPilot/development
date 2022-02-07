package main

import (
	"flag"
	"fmt"
	"net"

	"development/mxshop_srvs/user_srv/handler"
	"development/mxshop_srvs/user_srv/proto"

	"google.golang.org/grpc"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	fmt.Println(*ip, *port)

	// 注册服务
	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handler.UserServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	// 启动服务
	err = g.Serve(listen)
	if err != nil {
		panic(err)
	}
}
