package main

import (
	"net"
	"net/rpc"

	"development/OldPackageTest/new_helloworld/handler"
	"development/OldPackageTest/new_helloworld/server_proxy"
)

func main() {
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	err = server_proxy.RegisterHelloService(&handler.HelloService{})
	if err != nil {
		panic(err)
	}

	accept, err := listen.Accept()
	if err != nil {
		panic(err)
	}
	rpc.ServeConn(accept)
}
