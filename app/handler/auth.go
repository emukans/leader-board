package handler

import (
	"net/http"
	"strings"
)

func Auth(writer http.ResponseWriter, request *http.Request) {
	authHeader := request.Header.Get("Authorisation")
	splitToken := strings.Split(authHeader, " ")

	if len(splitToken) != 2 || splitToken[1] != "123" {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
