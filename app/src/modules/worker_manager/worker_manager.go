// Package worker_manager provides functionality to manage a pool of workers
// for executing asynchronous jobs.
package worker_manager

import (
	"os"
	"strconv"

	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
)

// Worker represents an individual worker that can execute jobs.
type Worker struct {
	ID       string
	JobQueue chan func()
	IsFree   bool
}

// NewWorker creates and initializes a new worker with the specified ID.
func NewWorker(id string) *Worker {
	logger.GetLogger().Info("Creating worker with ID: " + id)
	return &Worker{
		ID:       id,
		JobQueue: make(chan func()),
		IsFree:   true,
	}
}

// Start starts the worker, enabling it to execute jobs from its job queue.
func (w *Worker) Start() {
	logger.GetLogger().Info("Starting worker with ID: " + w.ID)
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				w.IsFree = false
				job()
				w.IsFree = true
			}
		}
	}()
}

// AddJob adds a new job to the worker's job queue.
func (w *Worker) AddJob(job func()) {
	logger.GetLogger().Info("Adding job to worker with ID: " + w.ID)
	w.JobQueue <- job
}

// IsAvailable checks if the worker is available to accept new jobs.
func (w *Worker) IsAvailable() bool {
	return w.IsFree
}

// WorkerManager manages a pool of workers.
type WorkerManager struct {
	Workers []*Worker
}

// NewWorkerManager creates and initializes a new worker manager with the specified number of initial workers.
func NewWorkerManager(initialWorkers int) *WorkerManager {
	wm := &WorkerManager{}
	for i := 0; i < initialWorkers; i++ {
		wm.AddWorker(NewWorker(strconv.Itoa(i)))
	}
	return wm
}

// AddWorker adds a new worker to the worker manager.
func (wm *WorkerManager) AddWorker(worker *Worker) {
	wm.Workers = append(wm.Workers, worker)
	worker.Start()
}

// GetAvailableWorker retrieves an available worker from the worker manager.
func (wm *WorkerManager) GetAvailableWorker() *Worker {
	for _, worker := range wm.Workers {
		if worker.IsAvailable() {
			return worker
		}
	}
	return nil
}

// AddJob adds a new job to the worker manager.
func (wm *WorkerManager) AddJob(job func()) {
	go func() {
		for {
			for _, worker := range wm.Workers {
				if worker.IsAvailable() {
					worker.AddJob(job)
					return
				}
			}
		}
	}()
}

// Start starts all workers in the worker manager.
func (wm *WorkerManager) Start() {
	for _, worker := range wm.Workers {
		worker.Start()
	}
}

// GetWorkerStatus returns the status of all workers in the worker manager.
func (wm *WorkerManager) GetWorkerStatus() map[string]bool {
	status := make(map[string]bool)
	for _, worker := range wm.Workers {
		status[worker.ID] = worker.IsAvailable()
	}
	return status
}

var vm *WorkerManager = nil

// GetWorkerManager returns a singleton instance of the worker manager.
// It initializes the worker pool size based on the environment variable HTTP_WORKER_POOL_SIZE.
func GetWorkerManager() *WorkerManager {
	if vm == nil {
		var pool_size_str = os.Getenv("HTTP_WORKER_POOL_SIZE")
		pool_size, err := strconv.Atoi(pool_size_str)
		if err != nil {
			logger.GetLogger().Fatal("Error parsing worker pool size: " + err.Error())
		}
		vm = NewWorkerManager(pool_size)
	}
	return vm
}
