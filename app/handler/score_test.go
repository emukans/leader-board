package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"leader-board/app/middleware"
	"leader-board/app/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessInsert(test *testing.T) {
	db, error := sql.Open("sqlite3", model.DBPath)
	if error != nil {
		test.Fatal(error)
	}
	model.DeleteScores(db)

	score := &model.PlayerScore{Name: "Foo", Score: 10}
	payload, error := json.Marshal(score)
	if error != nil {
		test.Fatal(error)
	}
	request, error := http.NewRequest("POST", "/api/v1/player/score", bytes.NewBuffer(payload))
	if error != nil {
		test.Fatal(error)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Score)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusNoContent {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	addedScore := model.FindScoreByName("foo", db)
	if addedScore == nil {
		test.Error("user is not inserted")
	}

	model.DeleteScores(db)
}

func TestFailedAuth(test *testing.T) {
	score := &model.PlayerScore{Name: "Foo", Score: 10}
	payload, error := json.Marshal(score)
	if error != nil {
		test.Fatal(error)
	}
	request, error := http.NewRequest("POST", "/api/v1/player/score", bytes.NewBuffer(payload))
	if error != nil {
		test.Fatal(error)
	}

	request.Header.Add("Authorization", "Bearer 123-fail")
	requestRecorder := httptest.NewRecorder()
	handler := middleware.Auth(http.HandlerFunc(Score))

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusForbidden {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestWrongHTTPMethod(test *testing.T) {
	methodList := []string {"GET", "PATCH", "DELETE", "HEAD"}
	for _, method := range methodList {
		request, error := http.NewRequest(method, "/api/v1/player/score", nil)
		if error != nil {
			test.Fatal(error)
		}

		requestRecorder := httptest.NewRecorder()
		handler := middleware.AllowedMethod(http.HandlerFunc(Score), http.MethodPost)

		handler.ServeHTTP(requestRecorder, request)

		if status := requestRecorder.Code; status != http.StatusMethodNotAllowed {
			test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
		}
	}
}

func TestScoreIsLessThanExisting(test *testing.T) {
	db, error := sql.Open("sqlite3", model.DBPath)
	if error != nil {
		test.Fatal(error)
	}
	model.DeleteScores(db)
	defer model.DeleteScores(db)

	score := model.PlayerScore{Name: "Foo", Score: 10}
	score.Save(db)
	score = model.PlayerScore{Name: "Foo", Score: 5}
	payload, error := json.Marshal(&score)
	if error != nil {
		test.Fatal(error)
	}

	request, error := http.NewRequest("POST", "/api/v1/player/score", bytes.NewBuffer(payload))
	if error != nil {
		test.Fatal(error)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Score)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusNoContent {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	addedScore := model.FindScoreByName("Foo", db)
	if addedScore.Score != 10 {
		test.Error("user score should not be updated")
	}
}

func TestScoreIsGreaterThanExisting(test *testing.T) {
	db, error := sql.Open("sqlite3", model.DBPath)
	if error != nil {
		test.Fatal(error)
	}
	model.DeleteScores(db)
	defer model.DeleteScores(db)

	score := model.PlayerScore{Name: "Foo", Score: 10}
	score.Save(db)
	score = model.PlayerScore{Name: "Foo", Score: 15}
	payload, error := json.Marshal(&score)
	if error != nil {
		test.Fatal(error)
	}

	request, error := http.NewRequest("POST", "/api/v1/player/score", bytes.NewBuffer(payload))
	if error != nil {
		test.Fatal(error)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Score)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusNoContent {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	addedScore := model.FindScoreByName("Foo", db)
	if addedScore.Score != 15 {
		test.Error("user score should be updated")
	}
}
