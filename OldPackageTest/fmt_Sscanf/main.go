package main

import "fmt"

func main() {
	header := "Bearer     111111111  \n Bearer 22"
	var rawJWT string

	n, err := fmt.Sscanf(header, "Bearer %s", &rawJWT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)

	fmt.Printf("%#v\n", rawJWT) // "111111111"
}
