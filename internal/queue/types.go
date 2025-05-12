package queue

import (
	"time"

	"github.com/henrywhitaker3/shorturl/internal/uuid"
)

var (
	Create Queue = "create"
	Click  Queue = "click"

	CreateTask Task = "create"
	ClickTask  Task = "click"
)

func mapTaskToQueue(task Task) Queue {
	switch task {
	case CreateTask:
		return Create
	case ClickTask:
		return Click
	default:
		return DefaultQueue
	}
}

type CreateJob struct {
	ID     uuid.UUID `json:"id"`
	Url    string    `json:"url"`
	Domain string    `json:"domain"`
}

type ClickJob struct {
	ID   uuid.UUID `json:"id"`
	IP   string    `json:"ip"`
	Time time.Time `json:"time"`
}
