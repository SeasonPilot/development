package main

import (
	"context"
	"fmt"
	"log"

	"development/OldPackageTest/loadbalance_test/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//grpc lb
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/user_srv?wait=14s&tag=primary",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 调用接口
	userClient := proto.NewUserClient(conn)
	for i := 0; i < 100; i++ {
		_, err = userClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    1,
			PSize: 5,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
