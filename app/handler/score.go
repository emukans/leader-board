package handler

import (
	"encoding/json"
	"leader-board/app/model"
	"net/http"
)


func Score(writer http.ResponseWriter, request *http.Request) {
	var payload model.PlayerScore
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		HandleInternalErr(err, writer)
		return
	}

	model.PlayerScore{
		Name:  payload.Name,
		Score: payload.Score,
	}.Save()

	writer.WriteHeader(http.StatusNoContent)
}
