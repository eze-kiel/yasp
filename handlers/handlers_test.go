package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eze-kiel/yasp/handlers"
)

func TestHomePage(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HomePage)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUploadPage(t *testing.T) {

	req := httptest.NewRequest("GET", "/upload", nil)
	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UploadPage)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDownloadPage(t *testing.T) {

	req := httptest.NewRequest("GET", "/download", nil)
	nr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.DownloadPage)

	handler.ServeHTTP(nr, req)

	if status := nr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
