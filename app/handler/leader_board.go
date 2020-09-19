package handler

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"leader-board/app/model"
	"net/http"
	"strconv"
)


type playerScore struct {
	Name string `json:"name"`
	Score int `json:"score"`
	Rank int `json:"rank"`
}

type leaderBoardResponse struct {
	Results []playerScore `json:"results"`
	NextPage int `json:"next_page"`
}

func LeaderBoard(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(writer, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "../../db/leader_board.db")
	if err != nil {
		panic(err)
	}

	pageNumber := 1
	page := request.URL.Query().Get("page")
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			http.Error(writer, "Page number is not valid", http.StatusBadRequest)
			return
		}
	}
	limit := 10
	offset := (pageNumber - 1) * limit

	scoreList := model.FindAllScores(db, limit, offset)
	var response leaderBoardResponse

	for rank, score := range scoreList {
		response.Results = append(response.Results, playerScore{
			Name:  score.Name,
			Score: score.Score,
			Rank:  offset + rank + 1,
		})
	}

	scoreCount := model.FindScoreCount(db)

	if scoreCount > (offset + limit) {
		pageNumber += 1
	} else {
		pageNumber = 0
	}
	response.NextPage = pageNumber

	writer.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	writer.Write(payload)

	writer.WriteHeader(http.StatusOK)
}
