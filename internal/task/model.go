package task

import (
	"fmt"
	"time"
)

type TaskPriority int
type TaskStatus int

const (
	LowPriority TaskPriority = iota
	MediumPriority
	HighPriority
)

const (
	StatusPending TaskStatus = iota
	StatusInProgress
	StatusCompleted
)

var priorityValueToLabelMap = map[TaskPriority]string{
	LowPriority:    "low",
	MediumPriority: "medium",
	HighPriority:   "high",
}

var priorityLabelToValueMap = map[string]TaskPriority{
	"low":    LowPriority,
	"medium": MediumPriority,
	"high":   HighPriority,
}

func (p TaskPriority) String() string {
	return priorityValueToLabelMap[p]
}

func ParsePriority(label string) (TaskPriority, error) {
	priorityValue, ok := priorityLabelToValueMap[label]
	if !ok {
		return 0, fmt.Errorf("invalid priority label: %s", label)
	}

	return priorityValue, nil
}

type Task struct {
	Id          string
	Title       string
	Priority    TaskPriority
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     time.Time
	CompletedOn time.Time
}
