package main

import (
	"context"
	"fmt"

	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var goodsClient proto.GoodsClient

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
