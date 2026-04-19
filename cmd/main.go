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

	var priority string
	var dueDate string
	var title string
	var status string

	tasks := make(map[string]task.Task, 0)

	flag.StringVar(&title, "title", "", "Title of the task being added")
	flag.StringVar(&priority, "priority", "low", "Priority level for the task being added")
	flag.StringVar(&dueDate, "due", "", "Due date for the task being added. If not provided, the due date will be marked as tomorrow EOD.")
	flag.StringVar(&status, "status", "pending", "Status of the task being added")

	flag.Parse()

	if title == "" {
		fmt.Println("Error: Please provide a task title")
		os.Exit(1)
	}

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
		Title:       title,
		Priority:    parsedPriority,
		Status:      task.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     parsedDueDate,
		CompletedOn: time.Time{},
	}

	fp, err := os.OpenFile("tasks.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error: Failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	tasks[newTask.Id] = newTask

	encoder := json.NewEncoder(fp)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Printf("Error: Failed to encode task: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Task added successfully")
}
