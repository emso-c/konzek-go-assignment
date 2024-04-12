// Package routers provides functions for registering HTTP routers and handlers for various endpoints.
//
// Endpoints:
// GET /tasks - Retrieves a list of tasks from the database based on pagination parameters.
// GET /task/{id} - Retrieves a task from the database based on the provided ID.
// POST /task - Creates a new task in the database.
// PUT /task/{id} - Updates an existing task in the database based on the provided ID.
// DELETE /task/{id} - Deletes a task from the database based on the provided ID.
//
// Usage:
// Use the RegisterTasksRouter function to register the tasks router with the provided Gorilla Mux router.
//
// Example:
// RegisterTasksRouter(router)
package routers

import (
	"net/http"
	"sync"

	"github.com/emso-c/konzek-go-assignment/src/api/controllers"
	"github.com/emso-c/konzek-go-assignment/src/database"
	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
	"github.com/emso-c/konzek-go-assignment/src/modules/worker_manager"
	"github.com/gorilla/mux"
)

// RegisterTasksRouter registers the routes related to tasks management.
func RegisterTasksRouter(router *mux.Router) {
	logger := logger.GetLogger()
	tc := controllers.NewTaskController()
	db := database.GetDatabase()

	taskRouter := router.PathPrefix("/").Subrouter()
	taskRouter.HandleFunc("/tasks", enqueueJob(tc.GetTasks(db))).Methods("GET")
	taskRouter.HandleFunc("/task/{id}", enqueueJob(tc.GetTask(db))).Methods("GET")
	taskRouter.HandleFunc("/task", enqueueJob(tc.CreateTask(db))).Methods("POST")
	taskRouter.HandleFunc("/task/{id}", enqueueJob(tc.UpdateTask(db))).Methods("PUT")
	taskRouter.HandleFunc("/task/{id}", enqueueJob(tc.DeleteTask(db))).Methods("DELETE")

	logger.Info("Tasks router registered")
}

// enqueueJob is a middleware function that enqueues the incoming HTTP handler function as a job to be processed by a worker.
// This allows handling requests concurrently while maintaining order.
func enqueueJob(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)

		worker_manager.GetWorkerManager().AddJob(func() {
			defer wg.Done()
			handlerFunc(w, r)
		})

		wg.Wait() // Wait until all jobs are done
	}
}
