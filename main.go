package main

import (
	"fmt"
	"leader-board/app/handler"
	"leader-board/app/middleware"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/player/score", middleware.AllowedMethod(middleware.Auth(http.HandlerFunc(handler.Score)), http.MethodPost))
	mux.Handle("/api/v1/leader-board", middleware.AllowedMethod(middleware.Auth(http.HandlerFunc(handler.LeaderBoard)), http.MethodGet))

	fmt.Println("A server is listening on port 8000")
	http.ListenAndServe(":8000", mux)
}

