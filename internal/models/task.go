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

var statusSlugs = map[Status]string{
	Todo:       "todo",
	InProgress: "in_progress",
	Done:       "done",
}

func (status Status) String() string {
	if val, ok := statusLabels[status]; ok {
		return val
	}
	return "Unknown"
}

func (status Status) IsValid() bool {
	_, ok := statusLabels[status]
	return ok
}

func (status Status) Slug() string {
	if val, ok := statusSlugs[status]; ok {
		return val
	}
	return "unknown"
}

func ParseStatusSlug(slug string) (Status, bool) {
	for status, slugValue := range statusSlugs {
		if slugValue == slug {
			return status, true
		}
	}
	return 0, false
}

type Task struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
