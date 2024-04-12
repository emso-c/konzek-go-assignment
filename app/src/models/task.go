package models

import (
	"time"
)

type Task struct {
	Id          uint // auto-increment by default
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time // server default is: TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	UpdatedAt   time.Time // server default is: TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
}

type CreateTaskRequest struct {
	Title       string
	Description string
	Status      string
}
