package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"leader-board/app/model"
	"net/http"
	"strconv"
	"time"
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

const (
	MonthlyLeaderBoardPeriod = "monthly"
	AllTimeLeaderBoardPeriod = "all-time"
)


func LeaderBoard(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(writer, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "../../db/leader_board.db")
	if err != nil {
		panic(err)
	}

	pageNumber, err := parsePageNumber(request)
	if err != nil {
		http.Error(writer, "Page number is not valid", http.StatusBadRequest)
		return
	}

	var period string
	period, err = parsePeriod(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var periodFrom time.Time
	if period == MonthlyLeaderBoardPeriod {
		periodFrom = time.Now().AddDate(0, -1, 0)
	}

	limit := 10
	offset := (pageNumber - 1) * limit

	scoreList := model.FindAllScores(db, limit, offset, periodFrom)
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

func parsePageNumber(request *http.Request) (int, error) {
	pageNumber := 1
	var err error

	page := request.URL.Query().Get("page")
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			return pageNumber, err
		}
	}
	return pageNumber, nil
}

func parsePeriod(request *http.Request) (string, error) {
	periodParameter := request.URL.Query().Get("period")

	if periodParameter == "" {
		return "",  nil
	}

	switch periodParameter {
	case AllTimeLeaderBoardPeriod, MonthlyLeaderBoardPeriod:
		return periodParameter, nil
	}

	return "", errors.New(fmt.Sprintf("Invalid period query parameter. Period should be either %s or %s", AllTimeLeaderBoardPeriod, MonthlyLeaderBoardPeriod))
}
