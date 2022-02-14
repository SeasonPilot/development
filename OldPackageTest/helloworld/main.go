package main

import (
	"fmt"

	proto1 "development/OldPackageTest/helloworld/proto"

	"google.golang.org/protobuf/proto"
)

func main() {
	req := proto1.HelloReq{
		Name: "kk",
	}

	rsp, _ := proto.Marshal(&req)
	fmt.Println(string(rsp))
}
