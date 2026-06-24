package models

import "time"

type Status int

const (
	Todo       Status = 1
	InProgress Status = 2
	Done       Status = 3
)

var statusLabels = map[Status]string{
	Todo:       "Todo",
	InProgress: "In Progress",
	Done:       "Done",
}

func (status Status) String() string {
	if val, ok := statusLabels[status]; ok {
		return val
	}
	return "Unknown"
}

type Task struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
