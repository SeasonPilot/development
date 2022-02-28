package main

import (
	"context"
	"fmt"
	"sync"

	"mxshop-srvs/inventory_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var inventoryClient proto.InventoryClient
var conn *grpc.ClientConn

func TestSetInv(goodsID, num int32) {
	_, err := inventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsID: goodsID,
		Num:     num,
	},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestDetail(goodsID int32) {
	rsp, err := inventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsID: goodsID,
	},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell(wg *sync.WaitGroup) {
	/*
		1. 第一件扣减成功： 第二件： 1. 没有库存信息 2. 库存不足
		2. 两件都扣减成功
	*/
	defer wg.Done()
	_, err := inventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsID: 421,
				Num:     1,
			},
			//{
			//	GoodsID: 422,
			//	Num:     20,
			//},
		},
	},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestReback() {
	_, err := inventoryClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsID: 421,
				Num:     10,
			},
			{
				GoodsID: 422,
				Num:     20,
			},
		},
	},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func Init() {
	var err error
	conn, err = grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	inventoryClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()

	var i int32
	//for i = 421; i < 840; i++ {
	//	TestSetInv(i, 100)
	//}

	//TestDetail(421)

	var wg sync.WaitGroup
	wg.Add(20)
	for i = 0; i < 20; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
	//TestReback()
	defer conn.Close()
}
