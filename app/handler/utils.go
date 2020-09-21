package handler

import (
	"log"
	"net/http"
)

func handleInternalErr(err error, writer http.ResponseWriter) {
	http.Error(writer, "Oops...something went wrong", http.StatusInternalServerError)
	log.Println(err)
}
