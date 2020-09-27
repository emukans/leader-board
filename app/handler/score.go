package handler

import (
	"encoding/json"
	"leader-board/app/model"
	"net/http"
)


func Score(writer http.ResponseWriter, request *http.Request) {
	var payload model.PlayerScore
	ErrWriter{writer: writer}.Then(func(self ErrWriter) ErrWriter {
		self.err = json.NewDecoder(request.Body).Decode(&payload)

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		model.PlayerScore{
			Name:  payload.Name,
			Score: payload.Score,
		}.Save()

		writer.WriteHeader(http.StatusNoContent)

		return self
	}).MaybeHandleInternalError()
}
