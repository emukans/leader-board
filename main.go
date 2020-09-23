package main

import (
	"flag"
	"fmt"
	"leader-board/app/handler"
	"leader-board/app/middleware"
	"leader-board/db"
	"net/http"
)


func main() {
	shouldSeedDB := flag.Bool("seed", false, "Seed database with sample scores")
	flag.Parse()

	if *shouldSeedDB {
		db.Seed()
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/player/score", middleware.AllowedMethod(middleware.Auth(http.HandlerFunc(handler.Score)), http.MethodPost))
	mux.Handle("/api/v1/leader-board", middleware.AllowedMethod(middleware.Auth(http.HandlerFunc(handler.LeaderBoard)), http.MethodGet))

	fmt.Println("A server is listening on port 8000")
	http.ListenAndServe(":8000", mux)
}

