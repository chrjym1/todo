package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Task represents a task in the database.
type Task struct {
	ID       int       `json:"id"`
	TaskName string    `json:"task_name"`
	Status   bool      `json:"status"`
	Date     time.Time `json:"date"`
}

// Add a new task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO task (task_name, status, date) VALUES (?, ?, ?)")
	checkErr(err)
	defer stmt.Close()

	task.Date = time.Now()
	_, err = stmt.Exec(task.TaskName, task.Status, task.Date)
	checkErr(err)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Task added successfully!")
}

// View all tasks
func viewTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task_name, status, date FROM task")
	checkErr(err)
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.TaskName, &task.Status, &task.Date)
		checkErr(err)
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Update a task's status
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE task SET status = ? WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(task.Status, task.ID)
	checkErr(err)

	fmt.Fprintln(w, "Task updated successfully!")
}

// Delete a task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM task WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)

	fmt.Fprintln(w, "Task deleted successfully!")
}
