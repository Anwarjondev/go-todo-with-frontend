package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Anwarjondev/go-todo-with-frontend/config"
	"github.com/Anwarjondev/go-todo-with-frontend/models"
	"github.com/gorilla/mux"
)

var view = template.Must(template.ParseFiles("./views/index.html"))

// Show displays all todos
func Show(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, item, completed FROM todos")
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.Id, &todo.Item, &todo.Completed); err != nil {
			log.Printf("Error scanning todo: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := models.View{
		Todos: todos,
	}

	if err := view.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Add creates a new todo
func Add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	item := r.FormValue("item")
	if item == "" {
		http.Error(w, "Item is required", http.StatusBadRequest)
		return
	}

	// For now, we'll use a default user_id of 1
	_, err := config.DB.Exec("INSERT INTO todos (item, user_id) VALUES ($1, $2)", item, 1)
	if err != nil {
		log.Printf("Error inserting todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete removes a todo
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Complete marks a todo as completed
func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec("UPDATE todos SET completed = true WHERE id = $1", id)
	if err != nil {
		log.Printf("Error completing todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
