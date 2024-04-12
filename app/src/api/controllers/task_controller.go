// Package controllers provides HTTP request handlers for managing tasks in the system.
// It includes methods for retrieving tasks, creating new tasks, updating existing tasks,
// and deleting tasks.
//
// Usage:
// Use the TaskController type to create a new instance of the controller.
// Use the GetTasks, GetTask, CreateTask, UpdateTask, and DeleteTask methods to handle HTTP requests.
//
// Example:
// tc := NewTaskController()
// http.HandleFunc("/api/tasks", tc.GetTasks(db))
// http.HandleFunc("/api/task/{id}", tc.GetTask(db))
// http.HandleFunc("/api/tasks", tc.CreateTask(db))
// http.HandleFunc("/api/tasks", tc.UpdateTask(db))
// http.HandleFunc("/api/task/{id}", tc.DeleteTask(db))
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emso-c/konzek-go-assignment/src/api/responses"
	"github.com/emso-c/konzek-go-assignment/src/models"
	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
	"github.com/gorilla/mux"
)

// TaskController represents the controller for handling task-related HTTP requests.
type TaskController struct{}

// NewTaskController creates a new instance of the TaskController.
func NewTaskController() *TaskController {
	return &TaskController{}
}

// GetTasks retrieves a list of tasks from the database based on pagination parameters.
// HTTP GET http://localhost:8080/api/tasks
func (tc *TaskController) GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = logger.GetLogger()
		logger.Info("GetTasks")

		// Parse pagination parameters from the query string
		pageStr := r.URL.Query().Get("page")
		sizeStr := r.URL.Query().Get("size")

		// Set default values if not provided
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = 10 // Default page size
		}

		// Calculate offset
		offset := (page - 1) * size

		var tasks []models.Task

		rows, err := db.Query(fmt.Sprintf("SELECT * FROM tasks ORDER BY id LIMIT %d OFFSET %d", size, offset))
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Error getting tasks from database")
			logger.Error("Error getting tasks from database: " + err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var task models.Task
			err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
			if err != nil {
				responses.Error(w, http.StatusInternalServerError, "Error scanning tasks from database")
				logger.Error("Error scanning tasks from database: " + err.Error())
				return
			}
			tasks = append(tasks, task)
		}

		logger.Info("Tasks retrieved successfully from database")

		if tasks == nil {
			responses.JSON(w, http.StatusNoContent, tasks)
			return
		}

		responses.JSON(w, http.StatusOK, tasks)
	}
}

// GetTask retrieves a task by its ID from the database.
// HTTP GET http://localhost:8080/api/task/{id}
func (tc *TaskController) GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = logger.GetLogger()
		logger.Info("GetTask")
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			responses.Error(w, http.StatusBadRequest, "ID is required")
			logger.Error("ID is required")
			return
		}
		logger.Info("ID is:" + vars["id"])
		id := vars["id"]
		if _, err := strconv.Atoi(id); err != nil {
			responses.Error(w, http.StatusBadRequest, "ID must be numeric")
			logger.Error("ID must be numeric")
			return
		}

		var task models.Task

		err := db.QueryRow("SELECT * FROM tasks WHERE id = $1", id).Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Could not get task from database")
			logger.Error("Error getting task from database" + err.Error())
			return
		}

		logger.Info("Task retrieved successfully from database")

		responses.JSON(w, http.StatusOK, task)
	}
}

// CreateTask creates a new task in the database based on the provided request body.
// Example:
// HTTP POST http://localhost:8080/api/tasks
// Content-Type: application/json
//
//	{
//		"title": "Task 1",
//		"description": "Description of task 1",
//		"status": "pending"
//	}
func (tc *TaskController) CreateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = logger.GetLogger()
		logger.Info("CreateTask")

		var req models.CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Error decoding request body")
			logger.Error("Error decoding request body:" + err.Error())
			return
		}

		_, err = db.Exec(fmt.Sprintf("INSERT INTO tasks (title, description, status) VALUES ('%s', '%s', '%s')", req.Title, req.Description, req.Status))
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Error inserting task into database")
			logger.Error("Error inserting task into database:" + err.Error())
			return
		}

		logger.Info("Task inserted successfully into database")

		responses.JSON(w, http.StatusCreated, nil)
	}
}

// UpdateTask updates an existing task in the database based on the provided request body.
// Example:
// HTTP PUT http://localhost:8080/api/tasks
// Content-Type: application/json
//
//	{
//		"id": 1,
//		"title": "Task 1",
//		"description": "Description of task 1",
//		"status": "completed"
//	}
func (tc *TaskController) UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = logger.GetLogger()
		logger.Info("UpdateTask")

		var task models.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Error decoding request body")
			logger.Error("Error decoding request body:" + err.Error())
			return
		}

		// _, err = db.Exec("UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5", task.Title, task.Description, task.Status, task.UpdatedAt, task.Id)
		_, err = db.Exec(fmt.Sprintf("UPDATE tasks SET title = '%s', description = '%s', status = '%s', updated_at = '%s' WHERE id = %d", task.Title, task.Description, task.Status, task.UpdatedAt, task.Id))

		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Error updating task in database")
			logger.Error("Error updating task in database:" + err.Error())
			return
		}

		logger.Info("Task updated successfully in database")

		responses.JSON(w, http.StatusOK, task)
	}
}

// DeleteTask deletes a task from the database based on its ID.
// Example:
// HTTP DELETE http://localhost:8080/api/task/{id}
func (tc *TaskController) DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var logger = logger.GetLogger()
		logger.Info("DeleteTask")
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			responses.Error(w, http.StatusBadRequest, "ID is required")
			logger.Error("ID is required")
			return
		}
		id := vars["id"]
		if _, err := strconv.Atoi(id); err != nil {
			responses.Error(w, http.StatusBadRequest, "ID must be numeric")
			logger.Error("ID must be numeric")
			return
		}

		_, err := db.Exec(fmt.Sprintf("DELETE FROM tasks WHERE id = %s", id))
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, "Error deleting task from database")
			logger.Error("Error deleting task from database:" + err.Error())
			return
		}

		logger.Info("Task deleted successfully from database")

		responses.JSON(w, http.StatusOK, nil)
	}
}
