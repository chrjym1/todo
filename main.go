package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

// Initialize the database connection
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// Check for errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Main function to set up HTTP routes
func main() {
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/add-task", addTaskHandler)
	mux.HandleFunc("/view-tasks", viewTasksHandler)
	mux.HandleFunc("/update-task", updateTaskHandler)
	mux.HandleFunc("/delete-task", deleteTaskHandler)

	//Server static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Add CORS support
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	corsAllowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	corsAllowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Use specific origins in production for better security

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsAllowedHeaders, corsAllowedMethods, corsAllowedOrigins)(mux)))
}
