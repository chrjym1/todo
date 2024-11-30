package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dixonwille/wmenu"
	_ "github.com/mattn/go-sqlite3"
)

type task struct {
	id        int
	task_name string
	status    bool
	date      *time.Time
}

func handleFunc(choice *string, db *sql.DB, opts []wmenu.Opt) {

	switch opts[0].Value {

	case 0:
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter new task: ")
		newTask, _ := reader.ReadString('\n')
		if newTask != "\n" {
			newTask = strings.TrimSuffix(newTask, "\n")
		}

		taskDate := time.Now()

		task := task{
			task_name: newTask,
			status:    false,
			date:      &taskDate,
		}

		addTask(db, task)

	case 1:
		tasks := showTasks(db)
		fmt.Printf("\n%-5s %-20s %-10s %-20s\n", "ID", "Task", "Finished", "Date       Time")
		fmt.Println(strings.Repeat("-", 70))

		for _, myTask := range tasks {
			var formattedDate string
			if myTask.date != nil {
				formattedDate = myTask.date.Format("2006-01-02 15:04:05")
			} else {
				formattedDate = "NULL"
			}
			fmt.Printf("%-5d %-20s %-10t %-20s\n", myTask.id, myTask.task_name, myTask.status, formattedDate)
		}

	case 2:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter task id to change status: ")
		updateID, _ := reader.ReadString('\n')

		currentTask := getTaskByID(db, updateID)

		affected := updateTask(db, currentTask)
		if affected == 1 {
			fmt.Println("One row affected")
		}

	case 3:
		fmt.Print("Enter the task ID to delete: ")
		var id int
		fmt.Scan(&id)

		deleteTask(db, id)

	case 4:
		fmt.Print("Bye.")
		*choice = "n"
	}

}

func addTask(db *sql.DB, task task) {
	stmt, err := db.Prepare("INSERT INTO TASK (task_name, status) VALUES (?,?)")
	checkErr(err)

	defer stmt.Close()

	_, err = stmt.Exec(task.task_name, task.status)
	checkErr(err)
}

func showTasks(db *sql.DB) []task {
	rows, err := db.Query("SELECT id, task_name, status, date FROM task")
	checkErr(err)

	defer rows.Close()

	tasks := make([]task, 0)

	for rows.Next() {
		var ourTask task
		var date sql.NullTime // Use sql.NullTime to handle NULL values
		err = rows.Scan(&ourTask.id, &ourTask.task_name, &ourTask.status, &date)
		checkErr(err)

		// Set the date field only if date is valid
		if date.Valid {
			// Convert the UTC time to your local timezone
			localTime := date.Time.In(time.Local) // or use a specific timezone with time.LoadLocation()
			ourTask.date = &localTime
		}
		tasks = append(tasks, ourTask)
	}
	return tasks
}

func getTaskByID(db *sql.DB, ourID string) task {
	rows, _ := db.Query("SELECT id, task_name, status FROM task WHERE id = '" + ourID + "'")
	defer rows.Close()

	ourTask := task{}

	for rows.Next() {
		rows.Scan(&ourTask.id, &ourTask.task_name, &ourTask.status)
	}

	return ourTask
}

func updateTask(db *sql.DB, currentTask task) int64 {
	stmt, err := db.Prepare("UPDATE task SET status = ? WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	if !currentTask.status {
		currentTask.status = true
	} else {
		currentTask.status = false
	}

	res, err := stmt.Exec(currentTask.status, currentTask.id)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}

func deleteTask(db *sql.DB, id int) {
	// Prepare the DELETE statement
	stmt, err := db.Prepare("DELETE FROM task WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	// Execute the DELETE statement
	result, err := stmt.Exec(id)
	checkErr(err)

	// Optionally, check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	checkErr(err)

	if rowsAffected == 0 {
		fmt.Printf("No task found with id %d\n", id)
	} else {
		fmt.Printf("Task with id %d deleted successfully.\n", id)
	}
}
