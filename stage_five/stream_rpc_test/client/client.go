package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"development/stage_five/stream_rpc_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	//服务端流模式
	client := proto.NewGreeterClient(conn)
	getStream, err := client.GetStream(context.Background(), &proto.StreamReqData{Data: "kk"})
	if err != nil {
		panic(err)
	}
	for {
		recv, err := getStream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(recv.GetData())
	}

	//客户端流模式
	putStream, err := client.PutStream(context.Background())
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		i++
		if i > 10 {
			break
		}
		err = putStream.Send(&proto.StreamReqData{
			Data: fmt.Sprintf("哈哈哈 %d\n", i),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
	// 双向流模式
	as, err := client.AllStream(context.Background())
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			_ = as.Send(&proto.StreamReqData{
				Data: "我是KK",
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
}
