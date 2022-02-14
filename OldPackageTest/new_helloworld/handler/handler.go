package handler

//名称冲突的问题
const HelloServiceName = "handler/HelloService"

type HelloService struct{}

// Hello 方法接受者为指针类型时，只能由 指针 类型变量调用
func (s *HelloService) Hello(request string, reply *string) error {
	//返回值是通过修改reply的值
	*reply = "hello, " + request
	return nil
}
