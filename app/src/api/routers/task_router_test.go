package routers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emso-c/konzek-go-assignment/config"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestEnqueueJob(t *testing.T) {
	os.Setenv("LOGGER_DISABLED", "true")
	os.Setenv("HTTP_WORKER_POOL_SIZE", "2")
	// Create a mock HTTP handler function
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the mock handler with enqueueJob
	wrappedHandler := enqueueJob(mockHandler)

	// Create a request to pass to the wrapped handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the wrapped handler with the mock request and response recorder
	wrappedHandler(rr, req)

	// Check if the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRegisterTasksRouter(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	cErr := config.LoadEnv("../../../config.toml")
	if cErr != nil {
		t.Fatal(cErr)
	}
	os.Setenv("LOGGER_DISABLED", "true")
	os.Setenv("HTTP_WORKER_POOL_SIZE", "2")

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Call RegisterTasksRouter with the router
	RegisterTasksRouter(router)

	// Create a new HTTP request for the "/tasks" endpoint
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check if the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
