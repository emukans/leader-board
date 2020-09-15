package main

import (
	"fmt"
	"leader-board/app/handler"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/auth", http.HandlerFunc(handler.Auth))
	
	fmt.Println("A server is listening on port 8000")
	http.ListenAndServe(":8000", mux)
}

