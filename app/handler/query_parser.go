package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	MonthlyLeaderBoardPeriod = "monthly"
	AllTimeLeaderBoardPeriod = "all-time"
)

func parseName(request *http.Request) string {
	name := request.URL.Query().Get("name")

	return name
}

func parsePageNumber(request *http.Request) (int, error) {
	pageNumber := 1
	var err error

	page := request.URL.Query().Get("page")
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if pageNumber <= 0 {
			err = errors.New("Page number should be greater than 0")
		}
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
