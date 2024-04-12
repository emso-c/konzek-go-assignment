package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/emso-c/konzek-go-assignment/config"
	"github.com/emso-c/konzek-go-assignment/src/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setup() {
	cErr := config.LoadEnv("../../../config.toml")
	if cErr != nil {
		log.Fatal(cErr)
	}
	os.Setenv("LOGGER_DISABLED", "true")
}

func TestCreateTask(t *testing.T) {
	setup()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new task
	task := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
	}
	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	expectedQuery := fmt.Sprintf("INSERT INTO tasks (title, description, status) VALUES ('%s', '%s', '%s')", task.Title, task.Description, task.Status)
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

	tc := NewTaskController()
	handler := http.HandlerFunc(tc.CreateTask(db))

	// Valid request
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(string(taskJSON)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Bad request
	req, err = http.NewRequest("POST", "/", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Valid params but invalid data
	req, err = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"title": "Test Task", "description": "Test Description", "status": 1, "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}`)))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mock.ExpectationsWereMet()
}

func TestUpdateTask(t *testing.T) {
	setup()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new task
	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedQuery := fmt.Sprintf("UPDATE tasks SET title = '%s', description = '%s', status = '%s', updated_at = '.*' WHERE id = %d", task.Title, task.Description, task.Status, task.Id)
	mock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, 1))

	tc := NewTaskController()
	handler := http.HandlerFunc(tc.UpdateTask(db))
	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	// Valid request
	req, err := http.NewRequest("PUT", "/1", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Bad request
	req, err = http.NewRequest("PUT", "/1", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Valid params but invalid data
	req, err = http.NewRequest("PUT", "/1", bytes.NewBuffer([]byte(`{"title": "Test Task", "description": "Test Description", "status": 1, "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}`)))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTask(t *testing.T) {
	setup()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new task
	task := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
		AddRow(task.Id, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)

	expectedQuery := "SELECT * FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(fmt.Sprintf("%d", task.Id)).
		WillReturnRows(rows)

	tc := NewTaskController()
	handler := http.HandlerFunc(tc.GetTask(db))

	// Valid request
	req, err := http.NewRequest("GET", fmt.Sprintf("/tasks/%d", task.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set mux vars
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", task.Id)})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Bad request
	req, err = http.NewRequest("GET", fmt.Sprintf("/tasks/%d", task.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// No need to set mux vars for this request

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTasks(t *testing.T) {
	setup()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new task
	task := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
		AddRow(task.Id, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)

	expectedQuery := "SELECT * FROM tasks"
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WillReturnRows(rows)

	tc := NewTaskController()
	handler := http.HandlerFunc(tc.GetTasks(db))

	// Valid request
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteTask(t *testing.T) {
	setup()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new task
	task := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedQuery := fmt.Sprintf("DELETE FROM tasks WHERE id = %d", task.Id)
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

	tc := NewTaskController()
	handler := http.HandlerFunc(tc.DeleteTask(db))

	// Valid request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", task.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set mux vars
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", task.Id)})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Bad request
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", task.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// No need to set mux vars for this request

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
