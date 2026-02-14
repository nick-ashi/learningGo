package dex

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func dex() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [args]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		// Reading/initializing the json
		data, err := os.ReadFile("tasks.json")
		var tasks []Task
		if err != nil {
			// Do nothing since it's already an empty slice
		} else {
			json.Unmarshal(data, &tasks)
			fmt.Println(tasks)
		}
		
		// Gen the task, append to task array
		var task Task = Task{
			ID: len(tasks),
			Text: string(os.Args[2]),
			Done: false,
		}
		tasks = append(tasks, task)

		// Convert tasks array to json
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			fmt.Println("Error marshaling task to JSON:", err)
			return
		}

		// Write the JSON byte slice to a file
		outputFile := "tasks.json"
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing task to file:", err)
			return
		}
		fmt.Printf("Added task: %+v\n", task)
		fmt.Println("Current Tasks: ", string(jsonData))

	case "list":
		// Reading/initializing the json
		data, err := os.ReadFile("tasks.json")
		var tasks []Task
		if err != nil {
			// Do nothing since it's already an empty slice
		} else {
			json.Unmarshal(data, &tasks)
			for i := 0; i < len(tasks); i++ {
				if (tasks[i].Done) {
					fmt.Print("ID:", i, " [X] ")
				} else {
					fmt.Print("ID:", i, " [ ] ")
				}
				fmt.Println(tasks[i].Text)
			}
			
		}
		
	case "done":
		// Reading/initializing the json
		id := os.Args[2]
		data, err := os.ReadFile("tasks.json")
		var tasks []Task
		if err != nil {
			// Do nothing since it's an empty slice
		} else {
			json.Unmarshal(data, &tasks)
		}

		var doneFlag bool
		for i := 0; i < len(tasks); i++ {
			if (strconv.Itoa(tasks[i].ID) == id) && !tasks[i].Done {
				tasks[i].Done = true
				doneFlag = true
			} else if (strconv.Itoa(tasks[i].ID) == id) && tasks[i].Done {
				tasks[i].Done = false
				doneFlag = false
			}
		}

		// Convert tasks array to json
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			fmt.Println("Error marshaling task to JSON:", err)
			return
		}

		// Write the JSON byte slice to a file
		outputFile := "tasks.json"
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing task to file:", err)
			return
		}

		if (doneFlag) {
			fmt.Println("Task", id + ": [X]")
			var ans string
			fmt.Print("Delete this task from the list? (y/n): ")
			fmt.Scan(&ans)
			if (ans == "y") {
				deleteTask()
			}
		} else {
			fmt.Println("Task", id + ": [ ]")
		}
		
	case "delall":
		var ans string
		fmt.Print("Are you sure you'd like to delete all tasks? (y/n): ")
		fmt.Scan(&ans)
		if ans == "y" {
			outputFile := "tasks.json"
			err := os.RemoveAll(outputFile)
			if err != nil {
				fmt.Println("Error removing file: ", err)
				return
			}
			fmt.Println("All tasks were deleted...")
		} else {
			fmt.Println("Nothing happened...")
		}

	case "rm":
		deleteTask()

	default:
		fmt.Println("Unknown command:", command)
	}
}

func deleteTask() {
	// Reading/initializing the json
	data, err := os.ReadFile("tasks.json")
	var tasks []Task
	if err != nil {
		// Do nothing since it's an empty slice
	} else {
		json.Unmarshal(data, &tasks)
	}
	
	// Gen the task, append to task array
	// ig the command line args are in scope still?
	id := os.Args[2]
	rmIndex := 0
	for i := 0; i < len(tasks); i++ {
		if strconv.Itoa(tasks[i].ID) == id {
			rmIndex = i
			break
		}
	}

	tasks = append(tasks[:rmIndex], tasks[rmIndex+1:]...)

	for i := 0; i < len(tasks); i++ {
		tasks[i].ID = i
	}

	// Convert tasks array to json
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error marshaling task to JSON:", err)
		return
	}

	// Write the JSON byte slice to a file
	outputFile := "tasks.json"
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing task to file:", err)
		return
	}

	fmt.Print("Task deleted")
}
