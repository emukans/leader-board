package handler

import (
	"database/sql"
	"encoding/json"
	"leader-board/app/model"
	"net/http"
)


type scorePayload struct {
	Name string `json:"name"`
	Score int `json:"score"`
}

func Score(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "../../db/leader_board.db")
	if err != nil {
		panic(err)
	}

	var payload scorePayload
	err = json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		panic(err)
	}

	model.PlayerScore{
		Name:  payload.Name,
		Score: payload.Score,
	}.Save(db)

	writer.WriteHeader(http.StatusNoContent)
}
