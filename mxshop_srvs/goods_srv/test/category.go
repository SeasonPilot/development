package main

import (
	"context"
	"fmt"

	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCategoryList() {
	list, err := goodsClient.GetAllCategoriesList(context.Background(), &emptypb.Empty{})
	if err != nil {
		return
	}

	fmt.Println(list.JsonData)
}

func TestGetSubCategory() {
	list, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
		//Level: 0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(list.SubCategories)
}

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
func main() {
	Init()
	defer conn.Close()

	//TestGetBrandList()
	//TestCategoryList()
	TestGetSubCategory()
}
