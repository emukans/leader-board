package middleware

import (
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorisation")
		splitToken := strings.Split(authHeader, "Bearer ")

		if len(splitToken) != 2 || splitToken[1] == "" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		} else if splitToken[1] != "123" {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		writer.WriteHeader(http.StatusOK)

		next.ServeHTTP(writer, request)
	})
}
