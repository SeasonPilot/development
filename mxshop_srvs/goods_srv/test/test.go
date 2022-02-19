package main

import (
	"context"
	"fmt"

	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var goodsClient proto.GoodsClient

func Init() {
	var err error
	conn, err = grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	goodsClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	resp, err := goodsClient.BrandList(context.TODO(), &proto.BrandFilterRequest{
		Pages:       1,
		PagePerNums: 3,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Total)
	for _, brand := range resp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()
	defer conn.Close()

	TestGetBrandList()
}
