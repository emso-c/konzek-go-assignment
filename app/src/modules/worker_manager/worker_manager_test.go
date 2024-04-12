package worker_manager

import (
	"os"
	"testing"
	"time"
)

func setup() {
	os.Setenv("HTTP_WORKER_POOL_SIZE", "2")
	os.Setenv("LOGGER_DISABLED", "true")
}

func TestWorker(t *testing.T) {
	setup()

	worker := NewWorker("1")
	if worker.ID != "1" {
		t.Errorf("Expected worker ID to be '1', got '%s'", worker.ID)
	}

	worker.Start()
	jobExecuted := false
	worker.AddJob(func() {
		jobExecuted = true
	})
	// Wait for the job to be executed
	time.Sleep(time.Millisecond * 100)
	if !jobExecuted {
		t.Error("Expected job to be executed by worker")
	}

	jobExecuted = false
	worker.AddJob(func() {
		jobExecuted = true
	})
	// Don't wait for the job to be executed
	if jobExecuted {
		t.Error("Expected job to be not done immediately by worker")
	}
}

func TestWorkerManager(t *testing.T) {
	setup()

	wm := NewWorkerManager(2)
	if len(wm.Workers) != 2 {
		t.Errorf("Expected number of workers to be 2, got %d", len(wm.Workers))
	}

	worker := wm.GetAvailableWorker()
	if worker == nil {
		t.Error("Expected an available worker, got nil")
	}

	jobExecuted := false
	wm.AddJob(func() {
		jobExecuted = true
	})

	time.Sleep(time.Millisecond * 100)
	if !jobExecuted {
		t.Error("Expected job to be executed by worker")
	}

	status := wm.GetWorkerStatus()
	for id, isAvailable := range status {
		if !isAvailable {
			t.Errorf("Expected worker %s to be available", id)
		}
	}
}
