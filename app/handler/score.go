package handler

import (
	"database/sql"
	"encoding/json"
	"leader-board/app/model"
	"net/http"
)


type scorePayload struct {
	name string
	score int
}

func Score(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)

	db, error := sql.Open("sqlite3", "../../db/leader_board.db")
	if error != nil {
		panic(error)
	}

	var payload scorePayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		panic(err)
	}

	model.PlayerScore{
		Name:  payload.name,
		Score: payload.score,
	}.Save(db)
}
