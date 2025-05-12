package queue

import "github.com/henrywhitaker3/go-template/internal/uuid"

var (
	Create Queue = "create"

	CreateTask Task = "create"
)

func mapTaskToQueue(task Task) Queue {
	switch task {
	case CreateTask:
		return Create
	default:
		return DefaultQueue
	}
}

type CreateJob struct {
	ID     uuid.UUID `json:"id"`
	Url    string    `json:"url"`
	Domain string    `json:"domain"`
}
