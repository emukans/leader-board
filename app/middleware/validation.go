package middleware

import (
	"net/http"
)

func AllowedMethod(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != method {
			http.Error(writer, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
