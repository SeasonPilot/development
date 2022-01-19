package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		return
	}

	var reply string
	err = client.Call("HelloService.Hello", "season", &reply)
	if err != nil {
		return
	}

	fmt.Println(reply)
}
