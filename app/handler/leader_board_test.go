package handler

import (
	"encoding/json"
	"leader-board/app/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestEmptyResponse(test *testing.T) {
	model.DeleteScores()

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
	limit := 5
	seedDb(limit)
	defer model.DeleteScores()

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
	limit := 15
	seedDb(limit)
	defer model.DeleteScores()

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
		if jsonBody.NextPage != page+1 && jsonBody.NextPage != 0 {
			test.Errorf("handler returned wrong body: next_page should be either %d or 0", page+1)
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
	limit := 15
	oldScoreCount := 8
	seedDb(limit)

	defer model.DeleteScores()

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

func TestExistingName(test *testing.T) {
	limit := 15
	seedDb(limit)

	defer model.DeleteScores()

	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	query := request.URL.Query()
	query.Add("name", "Dogge")
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

	expectedResultCount := 10
	if resultCount := len(jsonBody.AroundMe); resultCount != expectedResultCount {
		test.Errorf("handler returned wrong body: around_me should contain %d scores, but returned %d", expectedResultCount, resultCount)
	}

	if !isNameInScoreList(jsonBody.AroundMe, "Dogge") {
		test.Errorf("handler returned wrong body: Dogge is not in the around_me list")
	}
}

func TestNotExistingName(test *testing.T) {
	limit := 15
	seedDb(limit)

	defer model.DeleteScores()

	request, err := http.NewRequest("GET", "/api/v1/leader-board", nil)
	if err != nil {
		test.Fatal(err)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderBoard)

	query := request.URL.Query()
	query.Add("name", "NotExists")
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

	print(len(jsonBody.AroundMe))
	if jsonBody.AroundMe != nil {
		test.Errorf("handler returned wrong body: around_me should be nil")
	}
}

func seedDb(limit int) {
	model.DeleteScores()

	scoreList := []model.PlayerScore{
		{
			Name:      "Cat",
			Score:     31,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Dog",
			Score:     12,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Dogge",
			Score:     11,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Banye, The Omg Cat",
			Score:     31,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Puss in boots",
			Score:     63,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Beethoven",
			Score:     21,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Toto",
			Score:     155,
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			Name:      "Hund von baskerville",
			Score:     312,
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
		score.Save()
	}
}
