package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"leader-board/app/model"
	"net/http"
	"strconv"
	"time"
)


type leaderBoardResponse struct {
	Results []model.PlayerScore `json:"results"`
	NextPage int `json:"next_page"`
}

const (
	MonthlyLeaderBoardPeriod = "monthly"
	AllTimeLeaderBoardPeriod = "all-time"
)


func LeaderBoard(writer http.ResponseWriter, request *http.Request) {
	pageNumber, err := parsePageNumber(request)
	if err != nil {
		http.Error(writer, "Page number is not valid", http.StatusBadRequest)
		return
	}

	var periodFrom time.Time
	periodFrom, err = parsePeriod(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := model.FindPageLimit()
	if err != nil {
		HandleInternalErr(err, writer)
		return
	}
	offset := (pageNumber - 1) * limit

	scoreList, err := model.FindAllScores(limit, offset, periodFrom)
	if err != nil {
		HandleInternalErr(err, writer)
		return
	}
	response := leaderBoardResponse{Results: []model.PlayerScore{}, NextPage: 0}

	for rank, score := range scoreList {
		response.Results = append(response.Results, model.PlayerScore{
			Name:  score.Name,
			Score: score.Score,
			Rank:  offset + rank + 1,
		})
	}

	scoreCount, err := model.FindScoreCount()
	if err != nil {
		HandleInternalErr(err, writer)
		return
	}

	if scoreCount > (offset + limit) {
		response.NextPage = pageNumber + 1
	}

	writer.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response)
	if err != nil {
		HandleInternalErr(err, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(payload)
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

func parsePeriod(request *http.Request) (time.Time, error) {
	periodParameter := request.URL.Query().Get("period")

	var periodDate time.Time
	var err error

	if periodParameter == MonthlyLeaderBoardPeriod {
		periodDate = time.Now().AddDate(0, -1, 0)
	} else if periodParameter != "" && periodParameter != AllTimeLeaderBoardPeriod {
		err = errors.New(fmt.Sprintf("Invalid period query parameter. Period should be either %s or %s", AllTimeLeaderBoardPeriod, MonthlyLeaderBoardPeriod))
	}

	return periodDate, err
}
