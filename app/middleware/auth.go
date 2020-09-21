package middleware

import (
	"leader-board/app/model"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")

		if len(splitToken) != 2 || splitToken[1] == "" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		} else if authToken := model.FindAuthToken(); splitToken[1] != authToken {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
