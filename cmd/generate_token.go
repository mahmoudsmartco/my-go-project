package main

import (
	"app2_http_api_database/auth"
	"fmt"
)

func main() {
	token, err := auth.GenerateJWT("testuser")
	if err != nil {
		panic(err)
	}
	fmt.Println("Your JWT token:")
	fmt.Println(token)
}
