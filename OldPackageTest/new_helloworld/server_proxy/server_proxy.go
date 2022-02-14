package server_proxy

import (
	"net/rpc"

	"development/OldPackageTest/new_helloworld/handler"
)

/*
	如何做到解耦，做到只关心方法？  使用接口类型
	入参类型为什么 不用 空接口时，而用自定义接口时？ 自定义接口可以定义对象一定要有哪些方法
*/

type HelloServicer interface {
	Hello(request string, reply *string) error
}

// RegisterHelloService 如何做到解耦 - 我们关心的是函数 鸭子类型
// 注册服务逻辑
func RegisterHelloService(src HelloServicer) error {
	return rpc.RegisterName(handler.HelloServiceName, src)
}
