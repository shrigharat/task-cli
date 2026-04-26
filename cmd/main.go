package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-cli/internal/task"
	"time"
)

func exitWithError(message string) {
	fmt.Println("Error: ", message)
	os.Exit(1)
}

func addDefaultTasks(tasks []task.Task) []task.Task {
	tasks = append(tasks, task.CreateTask(1, "Buy groceries", task.MediumPriority))
	tasks = append(tasks, task.CreateTask(2, "Finish report", task.HighPriority))
	tasks = append(tasks, task.CreateTask(3, "Call Alice", task.LowPriority))

	return tasks
}

func validateRequiredArgument(argData string, argName string) string {
	processedData := strings.TrimSpace(argData)
	if processedData == "" {
		exitWithError(fmt.Sprintf("%s cannot be empty", argName))
	}

	return processedData
}

func validateTaskId(taskIdString string, tasksLength int) int {
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		exitWithError("Invalid task ID")
	}
	if taskId < 1 || taskId > tasksLength {
		exitWithError("Invalid task ID")
	}

	return taskId
}

func main() {
	args := os.Args[1:]
	tasks := addDefaultTasks(make([]task.Task, 0))

	if len(args) == 0 {
		fmt.Println("Error: Please provide an operation")
		os.Exit(1)
	}

	operation := args[0]

	switch operation {
	case "add":
		dataArgument := args[1]
		taskTitle := validateRequiredArgument(dataArgument, "Title")
		newTask := task.CreateTask(uint8(len(tasks)+1), taskTitle, task.LowPriority)
		tasks = append(tasks, newTask)
		fmt.Printf("Task \"%s\" added successfully\n", taskTitle)
	case "update":
		dataArgument := args[1]
		newTaskTitle := args[2]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), len(tasks))
		newTaskTitle = validateRequiredArgument(newTaskTitle, "New task title")
		tasks[taskId-1].Title = newTaskTitle
		fmt.Printf("Task \"%s\" updated successfully\n", newTaskTitle)
	case "delete":
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), len(tasks))
		tasks = append(tasks[0:taskId-1], tasks[taskId:]...)
		fmt.Printf("Task %d deleted successfully\n", taskId)
	case "list":
		filter := args[1]
		filterStatus, err := task.ParseStatus(filter)
		if err != nil {
			exitWithError(fmt.Sprintf("Invalid status filter: %s", filter))
		}
		filteredTasks := make([]task.Task, 0)
		switch filterStatus {
		case task.StatusPending:
			for _, currentTask := range tasks {
				if currentTask.Status == task.StatusPending {
					filteredTasks = append(filteredTasks, currentTask)
				}
			}
		case task.StatusInProgress:
			for _, currentTask := range tasks {
				if currentTask.Status == task.StatusInProgress {
					filteredTasks = append(filteredTasks, currentTask)
				}
			}
		case task.StatusCompleted:
			for _, currentTask := range tasks {
				if currentTask.Status == task.StatusCompleted {
					filteredTasks = append(filteredTasks, currentTask)
				}
			}
		default:
			filteredTasks = tasks
		}

		for index, currentTask := range tasks {
			fmt.Printf("%d. %s\n", index+1, currentTask.Title)
		}
	case "complete":
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), len(tasks))
		tasks[taskId-1].Status = task.StatusCompleted
		tasks[taskId-1].CompletedOn = time.Now()
		fmt.Printf("Task %d completed successfully\n", taskId)
	case "mark-in-progress":
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), len(tasks))
		tasks[taskId-1].Status = task.StatusInProgress
		tasks[taskId-1].UpdatedAt = time.Now()
		fmt.Printf("Task %d marked as in progress successfully\n", taskId)
	case "mark-pending":
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), len(tasks))
		tasks[taskId-1].Status = task.StatusPending
		tasks[taskId-1].UpdatedAt = time.Now()
		fmt.Printf("Task %d marked as pending successfully\n", taskId)
	default:
		fmt.Println("Error: Invalid operation")
		os.Exit(1)
	}
}
