package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Anwarjondev/go-todo-with-frontend/config"
	"github.com/Anwarjondev/go-todo-with-frontend/controllers"
	"github.com/gorilla/mux"
)

// @title Todo API
// @version 1.0
// @description This is a todo list API
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// API routes
	router.HandleFunc("/", controllers.Show).Methods("GET")
	router.HandleFunc("/add", controllers.Add).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}", controllers.Delete).Methods("GET")
	router.HandleFunc("/complete/{id:[0-9]+}", controllers.Complete).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
