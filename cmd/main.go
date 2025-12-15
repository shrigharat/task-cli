package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"task-cli/internal/task"
	"time"

	"github.com/google/uuid"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Error: Please provide a task title")
		os.Exit(1)
	}

	taskTitle := args[0]

	if taskTitle == "" {
		fmt.Println("Error: Please provide a task title")
		os.Exit(1)
	}

	flag.Parse()

	var priority string
	var dueDate string

	flag.StringVar(&priority, "priority", "low", "Priority level for the task being added")
	flag.StringVar(&dueDate, "due", "", "Due date for the task being added")

	var parsedDueDate time.Time
	parsedPriority, err := task.ParsePriority(priority)
	if err != nil {
		fmt.Printf("Error: Invalid priority level: %s\n", err)
		os.Exit(1)
	}

	if dueDate != "" {
		parsedDueDate, err = time.Parse("2006-01-02", dueDate)
		if err != nil {
			fmt.Println("Waring: Date provided is not the correct format (YYYY-MM-DD)")
		}
	}

	newTask := task.Task{
		Id:          uuid.New().String(),
		Title:       taskTitle,
		Priority:    parsedPriority,
		Status:      task.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     parsedDueDate,
		CompletedOn: time.Time{},
	}

	fp, err := os.Create("tasks.json")
	if err != nil {
		fmt.Printf("Error: Failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	encoder := json.NewEncoder(fp)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(newTask)
	if err != nil {
		fmt.Printf("Error: Failed to encode task: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Task added successfully")
}
