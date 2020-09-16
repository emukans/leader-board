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
	db, error := sql.Open("sqlite3", "../../db/leader_board.db")
	if error != nil {
		test.Fatal(error)
	}
	cleanTestUser(db)

	score := &model.PlayerScore{Name: "Foo", Score: 10}
	payload, error := json.Marshal(score)
	if error != nil {
		test.Fatal(error)
	}
	request, error := http.NewRequest("POST", "/api/v1/player/score", bytes.NewBuffer(payload))
	if error != nil {
		test.Fatal(error)
	}

	request.Header.Add("Authorisation", "Bearer 123")
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Score)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusNoContent {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	scoreList := model.FindByName("foo", db)
	if len(scoreList) != 1 {
		test.Error("user is not inserted")
	}

	cleanTestUser(db)
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

	request.Header.Add("Authorisation", "Bearer 123-fail")
	requestRecorder := httptest.NewRecorder()
	handler := middleware.Auth(http.HandlerFunc(Score))

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusForbidden {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func cleanTestUser(db *sql.DB) {
	scoreList := model.FindByName("foo", db)
	if len(scoreList) > 0 {
		for _, score := range scoreList {
			score.Delete(db)
		}
	}
}
