package handler

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"leader-board/app/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)


func TestEmptyResponse(test *testing.T) {
	db, err := sql.Open("sqlite3", model.DBPath)
	if err != nil {
		test.Fatal(err)
	}

	model.DeleteScores(db)

	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	jsonBody := leaderBoardResponse{}
	err = json.NewDecoder(requestRecorder.Body).Decode(&jsonBody)

	if err != nil {
		test.Fatal(err)
	}
	if len(jsonBody.Results) != 0 {
		test.Errorf("handler returned wrong body: results should be empty")
	}
	if jsonBody.NextPage != 0 {
		test.Errorf("handler returned wrong body: next_page should be 0")
	}
}

func TestOnePageSeededDb(test *testing.T) {
	db, err := sql.Open("sqlite3", model.DBPath)
	if err != nil {
		test.Fatal(err)
	}
	limit := 5
	seedDb(db, limit)
	defer model.DeleteScores(db)

	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	jsonBody := leaderBoardResponse{}
	err = json.NewDecoder(requestRecorder.Body).Decode(&jsonBody)

	if err != nil {
		test.Fatal(err)
	}
	if len(jsonBody.Results) != limit {
		test.Errorf("handler returned wrong body: results should contain %d scores", limit)
	}
	if jsonBody.NextPage != 0 {
		test.Errorf("handler returned wrong body: next_page should be 0")
	}
}

func TestMultiPageSeededDb(test *testing.T) {
	db, err := sql.Open("sqlite3", model.DBPath)
	if err != nil {
		test.Fatal(err)
	}
	limit := 15
	seedDb(db, limit)
	defer model.DeleteScores(db)

	page := 1

	for page != 0 {
		request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
		if err != nil {
			test.Fatal(err)
		}

		requestRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(LeaderBoard)

		query := request.URL.Query()
		query.Add("page", strconv.Itoa(page))
		request.URL.RawQuery = query.Encode()

		handler.ServeHTTP(requestRecorder, request)

		if status := requestRecorder.Code; status != http.StatusOK {
			test.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		jsonBody := leaderBoardResponse{}
		err = json.NewDecoder(requestRecorder.Body).Decode(&jsonBody)

		if err != nil {
			test.Fatal(err)
		}
		if jsonBody.NextPage != page + 1 && jsonBody.NextPage != 0 {
			test.Errorf("handler returned wrong body: next_page should be either %d or 0", page + 1)
		}

		page = jsonBody.NextPage
	}
}


func TestFailedPeriod(test *testing.T) {
	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	query := request.URL.Query()
	query.Add("period", "wrong-period")
	request.URL.RawQuery = query.Encode()

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusBadRequest {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestMonthlyPeriod(test *testing.T) {
	db, err := sql.Open("sqlite3", model.DBPath)
	if err != nil {
		test.Fatal(err)
	}
	limit := 15
	oldScoreCount := 8
	seedDb(db, limit)

	defer model.DeleteScores(db)

	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	query := request.URL.Query()
	query.Add("period", MonthlyLeaderBoardPeriod)
	request.URL.RawQuery = query.Encode()

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	jsonBody := leaderBoardResponse{}
	err = json.NewDecoder(requestRecorder.Body).Decode(&jsonBody)

	if err != nil {
		test.Fatal(err)
	}
	expectedResultCount := limit - oldScoreCount
	actualResultCount := len(jsonBody.Results)
	if actualResultCount != expectedResultCount {
		test.Errorf("handler returned wrong body: results should contain %d scores, but returned %d", expectedResultCount, actualResultCount)
	}
}


func seedDb(db *sql.DB, limit int) {
	model.DeleteScores(db)

	scoreList := []model.PlayerScore {
		{
			Name:  "Cat",
			Score: 1,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Dog",
			Score: 12,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Dogge",
			Score: 11,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Banye, The Omg Cat",
			Score: 31,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Puss in boots",
			Score: 63,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Beethoven",
			Score: 12,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Toto",
			Score: 155,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Hund von baskerville",
			Score: 312,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:  "Grumpy cat",
			Score: 39,
		},
		{
			Name:  "Lil Bub",
			Score: 93,
		},
		{
			Name:  "Maru, The master of boxes",
			Score: 221,
		},
		{
			Name:  "Garfield",
			Score: 43,
		},
		{
			Name:  "Hamilton, The Hipster Cat",
			Score: 33,
		},
		{
			Name:  "Waffles The Cat",
			Score: 54,
		},
		{
			Name:  "Kitty",
			Score: 1,
		}}

	for _, score := range scoreList[:limit] {
		score.Save(db)
	}
}

