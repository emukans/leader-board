package main

import (
	"fmt"
	"leader-board/app/handler"
	"leader-board/app/middleware"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/player/score", middleware.Auth(http.HandlerFunc(handler.Score)))
	mux.Handle("/api/v1/leader-board", middleware.Auth(http.HandlerFunc(handler.LeaderBoard)))

	fmt.Println("A server is listening on port 8000")
	http.ListenAndServe(":8000", mux)
}

