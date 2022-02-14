package main

import (
	"fmt"

	"development/OldPackageTest/new_helloworld/client_proxy"
)

func main() {
	client := client_proxy.NewHelloServiceClient("tcp", "localhost:1234")

	var reply string
	err := client.Hello("season", &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
