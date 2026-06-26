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
	Title       string    `json:"title" validate:"required,min=5,max=255"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Status      Status    `json:"status" validate:"required,oneof=1 2 3"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
