package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/publish", publishHandler)

	fmt.Println("server is running at localhost:3000")
	http.ListenAndServe(":3000", nil)
}
