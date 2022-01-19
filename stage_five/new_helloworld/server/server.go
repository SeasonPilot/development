package main

import (
	"development/stage_five/new_helloworld/handler"
	"net"
	"net/rpc"
)

func main() {
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	err = rpc.RegisterName(handler.HelloServiceName, &handler.HelloService{})
	if err != nil {
		panic(err)
	}

	accept, err := listen.Accept()
	if err != nil {
		panic(err)
	}
	rpc.ServeConn(accept)
}
