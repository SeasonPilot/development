package main

import (
	"net"
	"net/rpc"
)

type HelloService struct{}

func (s *HelloService) Hello(request string, reply *string) error {
	//返回值是通过修改reply的值
	*reply = "hello, " + request
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	err = rpc.Register(&HelloService{})
	if err != nil {
		return
	}

	accept, err := listen.Accept()
	if err != nil {
		return
	}
	rpc.ServeConn(accept)
}
