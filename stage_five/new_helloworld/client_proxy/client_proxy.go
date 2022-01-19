package client_proxy

import (
	"net/rpc"

	"development/stage_five/new_helloworld/handler"
)

type HelloServiceStub struct {
	*rpc.Client
}

// NewHelloServiceClient 这里要把 rpc.Dial 封装进来，让用户只传入 network, addr即可
func NewHelloServiceClient(network, addr string) *HelloServiceStub {
	client, err := rpc.Dial(network, addr)
	if err != nil {
		panic("connect error!" + err.Error())
	}
	return &HelloServiceStub{Client: client}
}

func (c *HelloServiceStub) Hello(args string, reply *string) error {
	return c.Call(handler.HelloServiceName+".Hello", args, reply)
}
