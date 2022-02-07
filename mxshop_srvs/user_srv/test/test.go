package main

import (
	"context"
	"fmt"

	"development/mxshop_srvs/user_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var userClient proto.UserClient

func Init() {
	var err error
	conn, err = grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

func TestUserServer_CreateUser() {
	for i := 10; i < 20; i++ {
		resp, err := userClient.CreateUser(context.TODO(), &proto.CreateUserInfo{
			Mobile:   fmt.Sprintf("187987678%d", i),
			Password: "admin123",
			NickName: fmt.Sprintf("bobby%d", i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Id, resp.Mobile)
	}
}

func TestUserServer_GetUserList() {
	resp, err := userClient.GetUserList(context.TODO(), &proto.PageInfo{
		Pn:    3,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range resp.Data {
		checkResp, err := userClient.CheckPassWord(context.TODO(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkResp.Success)
	}
}

func main() {
	Init()
	defer conn.Close()

	TestUserServer_CreateUser()
	TestUserServer_GetUserList()
}
