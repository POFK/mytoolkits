package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Todo represents a to-do item
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fs := http.FileServer(http.Dir("./todoapp/build"))
	http.Handle("/", fs)

	http.HandleFunc("/todos", handleTodo)
	http.HandleFunc("/todos/", handleTodo)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func getId(w http.ResponseWriter, r *http.Request) int {
	id, err := strconv.Atoi(r.URL.Path[len("/todos/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return -1
	}
	return id
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println(">>>", r.URL, r.Method, r.Body, "\n")
		handleTodoGet(w, r)
	case "POST":
		fmt.Println(">>>", r.URL, r.Method, r.Body, "\n")
		handleTodoPost(w, r)
	case "PUT":
		fmt.Println(">>>", r.URL, r.Method, r.Body, "\n")
		id := getId(w,r)
		handleTodoPut(w, r, id)
	case "DELETE":
		fmt.Println(">>>", r.URL, r.Method, r.Body, "\n")
		id := getId(w,r)
		handleTodoDelete(w, r, id)
	default:
		fmt.Println(">>>", r.URL, r.Method, r.Body, "\n")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTodoGet(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		http.Error(w, "Error fetching to-dos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Text)
		if err != nil {
			http.Error(w, "Error scanning to-do", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error fetching to-dos", http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(todos)

	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleTodoPost(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO todos (text) VALUES (?)", todo.Text)
	if err != nil {
		http.Error(w, "Error inserting to-do", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting ID of inserted to-do", http.StatusInternalServerError)
		return
	}

	todo.ID = int(id)

	js, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleTodoDelete(w http.ResponseWriter, r *http.Request, id int) {
	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting to-do", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleTodoPut(w http.ResponseWriter, r *http.Request, id int) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	todo.ID = id

	_, err = db.Exec("UPDATE todos SET text = ? WHERE id = ?", todo.Text, todo.ID)
	if err != nil {
		http.Error(w, "Error updating to-do", http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
