package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrectAuthToken(test *testing.T) {
	request, error := http.NewRequest("GET", "/auth", nil)
	if error != nil {
		test.Fatal(error)
	}

	request.Header.Add("Authorisation", "Bearer 123")
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Auth)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNoAuthToken(test *testing.T) {
	request, error := http.NewRequest("GET", "/auth", nil)
	if error != nil {
		test.Fatal(error)
	}

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Auth)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusForbidden {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	// Check the response body is what we expect.
	//expected := `Hello`
	//if requestRecorder.Body.String() != expected {
	//	test.Errorf("handler returned unexpected body: got %v want %v",
	//		requestRecorder.Body.String(), expected)
	//}
}

func TestInvalidAuthToken(test *testing.T) {
	request, error := http.NewRequest("GET", "/auth", nil)
	if error != nil {
		test.Fatal(error)
	}

	request.Header.Add("Authorisation", "Bearer 321")
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Auth)

	handler.ServeHTTP(requestRecorder, request)

	if status := requestRecorder.Code; status != http.StatusForbidden {
		test.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}
