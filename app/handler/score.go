package handler

import (
	"database/sql"
	"encoding/json"
	"leader-board/app/model"
	"net/http"
)


func Score(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("sqlite3", model.DBPath)
	if err != nil {
		handleInternalErr(err, writer)
		return
	}

	var payload model.PlayerScore
	err = json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		handleInternalErr(err, writer)
		return
	}

	model.PlayerScore{
		Name:  payload.Name,
		Score: payload.Score,
	}.Save(db)

	writer.WriteHeader(http.StatusNoContent)
}
