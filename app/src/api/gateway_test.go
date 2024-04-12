package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emso-c/konzek-go-assignment/config"
	"github.com/joho/godotenv"
)

func TestInit(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	cErr := config.LoadEnv("../../config.toml")
	if cErr != nil {
		t.Fatal(cErr)
	}
	os.Setenv("LOGGER_DISABLED", "true")
	os.Setenv("HTTP_RATE_LIMIT", "9999")

	// Call the Init function to initialize the router
	Init()

	// Create a new HTTP request for testing
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check if the status code is what we expect (200 OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Test if CORS middleware is applied
	req, err = http.NewRequest("OPTIONS", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Header().Get("Access-Control-Allow-Origin") == "GET" {
		t.Errorf("CORS headers not set correctly")
	}

	// Test if NotFoundMiddleware is set
	req, err = http.NewRequest("GET", "/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for 404: got %v want %v",
			status, http.StatusNotFound)
	}

	// Test if RateLimitMiddleware is set
	req, err = http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for rate limit: got %v want %v",
			status, http.StatusOK)
	}

	// Test if CSRFMiddleware is set
	req, err = http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for CSRF: got %v want %v",
			status, http.StatusOK)
	}
	if rr.Header().Get("X-CSRF-Token") == "" {
		t.Errorf("CSRF token not set in response header")
	}

	// Test if SQLInjectionMiddleware is set
	req, err = http.NewRequest("GET", "/api/tasks?page=1&size=1%20OR%201=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for SQL injection: got %v want %v",
			status, http.StatusBadRequest)
	}
}
