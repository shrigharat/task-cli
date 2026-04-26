package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"task-cli/internal/task"
	"text/tabwriter"
	"time"
)

func exitWithError(message string) {
	fmt.Println("Error: ", message)
	os.Exit(1)
}

func validateRequiredArgument(argData string, argName string) string {
	processedData := strings.TrimSpace(argData)
	if processedData == "" {
		exitWithError(fmt.Sprintf("%s cannot be empty", argName))
	}

	return processedData
}

func validateTaskId(taskIdString string, tasks []task.Task) int {
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		exitWithError("Invalid task ID")
	}
	if taskId <= 0 {
		exitWithError("Invalid task ID")
	}

	taskIndex := getTaskIndexById(tasks, taskId)
	if taskIndex == -1 {
		exitWithError("Invalid task ID")
	}

	return taskIndex
}

func checkMinimumArguments(args []string, minimumArguments int) {
	if len(args) < minimumArguments {
		exitWithError(fmt.Sprintf("Please provide at least %d arguments", minimumArguments))
	}
}

func writeTasksToFile(tasks []task.Task) {
	newFileData := make(map[string][]task.Task)
	newFileData["tasks"] = tasks

	jsonData, err := json.MarshalIndent(newFileData, "", "\t")
	if err != nil {
		exitWithError(fmt.Sprintf("Error marshalling tasks: %s", err))
	}
	err = os.WriteFile("tasks.json", jsonData, 0644)
	if err != nil {
		exitWithError(fmt.Sprintf("Error writing to tasks file: %s", err))
	}
	fmt.Println("Tasks saved successfully")
}

func getMaxTaskId(tasks []task.Task) int {
	maxId := float64(0)
	for _, t := range tasks {
		maxId = math.Max(float64(maxId), float64(t.Id))
	}
	return int(maxId)
}

func getTaskIndexById(tasks []task.Task, taskId int) int {
	for index, t := range tasks {
		if t.Id == taskId {
			return index
		}
	}
	return -1
}

func main() {
	args := os.Args[1:]
	fileBytes, err := os.ReadFile("tasks.json")
	if err != nil {
		if !os.IsNotExist(err) {
			exitWithError(fmt.Sprintf("Error reading from tasks file: %s", err))
		}
	}

	tasks := make([]task.Task, 0)

	if len(fileBytes) > 0 {
		existingFileData := make(map[string][]task.Task)
		err = json.Unmarshal(fileBytes, &existingFileData)
		if err != nil {
			exitWithError(fmt.Sprintf("Error unmarshalling tasks: %s", err))
		}
		tasks = existingFileData["tasks"]
	}

	if len(args) == 0 {
		fmt.Println("Error: Please provide an operation")
		os.Exit(1)
	}

	operation := args[0]

	switch operation {
	case "add":
		checkMinimumArguments(args, 2)
		dataArgument := args[1]
		newTaskId := getMaxTaskId(tasks) + 1
		taskTitle := validateRequiredArgument(dataArgument, "Title")
		newTask := task.CreateTask(newTaskId, taskTitle, task.LowPriority)
		tasks = append(tasks, newTask)
		fmt.Printf("Task \"%s\" added successfully\n", taskTitle)
	case "update":
		checkMinimumArguments(args, 3)
		dataArgument := args[1]
		newTaskTitle := args[2]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), tasks)
		newTaskTitle = validateRequiredArgument(newTaskTitle, "New task title")
		tasks[taskId].Title = newTaskTitle
		tasks[taskId].UpdatedAt = time.Now()
		fmt.Printf("Task %d updated successfully\n", tasks[taskId].Id)
	case "delete":
		checkMinimumArguments(args, 2)
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), tasks)
		originalTask := tasks[taskId]
		tasks = append(tasks[0:taskId], tasks[taskId+1:]...)
		fmt.Printf("Task %d deleted successfully\n", originalTask.Id)
	case "list":
		checkMinimumArguments(args, 1)
		filteredTasks := make([]task.Task, 0)
		if len(args) <= 1 {
			filteredTasks = tasks
		} else {
			filter := args[1]
			filterStatus, err := task.ParseStatus(filter)
			if err != nil {
				exitWithError(fmt.Sprintf("Invalid status filter: %s", filter))
			}
			switch filterStatus {
			case task.StatusTodo:
				for _, currentTask := range tasks {
					if currentTask.Status == task.StatusTodo {
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
		}

		if len(filteredTasks) == 0 {
			fmt.Println("There are no tasks matching the filter")
		} else {
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(tabWriter, "ID\tTitle\tStatus\tPriority\n")
			for _, currentTask := range filteredTasks {
				priorityLabel := task.GetPriorityLabel(currentTask.Priority)
				statusLabel := task.GetStatusLabel(currentTask.Status)
				fmt.Fprintf(tabWriter, "%d.\t%s\t%s\t%s\n", currentTask.Id, currentTask.Title, statusLabel, priorityLabel)
			}
			tabWriter.Flush()
		}
	case "complete":
		checkMinimumArguments(args, 2)
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), tasks)
		tasks[taskId].Status = task.StatusCompleted
		tasks[taskId].CompletedOn = time.Now()
		fmt.Printf("Task %d completed successfully\n", tasks[taskId].Id)
	case "mark-in-progress":
		checkMinimumArguments(args, 2)
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), tasks)
		tasks[taskId].Status = task.StatusInProgress
		tasks[taskId].UpdatedAt = time.Now()
		fmt.Printf("Task %d marked as in progress successfully\n", tasks[taskId].Id)
	case "mark-pending":
		checkMinimumArguments(args, 2)
		dataArgument := args[1]
		taskId := validateTaskId(validateRequiredArgument(dataArgument, "Task ID"), tasks)
		tasks[taskId].Status = task.StatusTodo
		tasks[taskId].UpdatedAt = time.Now()
		fmt.Printf("Task %d marked as pending successfully\n", tasks[taskId].Id)
	default:
		fmt.Println("Error: Invalid operation")
		os.Exit(1)
	}

	if operation != "list" {
		writeTasksToFile(tasks)
	}
}
