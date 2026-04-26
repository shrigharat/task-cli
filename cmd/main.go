package main

import (
	"fmt"
	"os"
	"task-cli/internal/task"
	"time"
)

func main() {
	args := os.Args[1:]
	tasks := make([]task.Task, 0)

	fmt.Println("Args: ", args)
	if len(args) == 0 {
		fmt.Println("Error: Please provide an operation")
		os.Exit(1)
	}

	operation := args[0]
	taskTitle := args[1]

	switch operation {
	case "add":
		newTask := task.Task{
			Id:          uint8(len(tasks) + 1),
			Title:       taskTitle,
			Priority:    task.LowPriority,
			Status:      task.StatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DueDate:     time.Time{},
			CompletedOn: time.Time{},
		}
		tasks = append(tasks, newTask)
		fmt.Printf("Task \"%s\" added successfully\n", taskTitle)
	case "list":
		fmt.Println("Listing tasks")
	case "complete":
		fmt.Println("Completing task: ", taskTitle)
	case "delete":
		fmt.Println("Deleting task: ", taskTitle)
	default:
		fmt.Println("Error: Invalid operation")
		os.Exit(1)
	}
}
