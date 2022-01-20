package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"development/stage_five/stream_rpc_test/proto"

	"google.golang.org/grpc"
)

type Server struct{}

// GetStream 服务端流模式
func (s *Server) GetStream(req *proto.StreamReqData, srv proto.Greeter_GetStreamServer) error {
	fmt.Println(req)
	i := 0
	for {
		i++
		if i > 10 {
			break
		}
		// send 数据
		err := srv.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%d\n", time.Now().Unix()),
		})
		if err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
	return nil
}

// PutStream 客户端流模式
func (s *Server) PutStream(srv proto.Greeter_PutStreamServer) error {
	for {
		recv, err := srv.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(recv.Data)
	}
	return nil
}
func (s *Server) AllStream(as proto.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			_ = as.Send(&proto.StreamResData{
				Data: "我是服务端",
			})
			time.Sleep(time.Second)
		}

	}()

	go func() {
		defer wg.Done()
		for {
			recv, _ := as.Recv()
			fmt.Println(recv.Data)
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()

	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	err = g.Serve(listener)
	if err != nil {
		panic(err)
	}
}
