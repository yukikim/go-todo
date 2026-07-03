package main

import (
	"fmt"
	"log"
	"os"

	"go-todo/internal/db"
	"go-todo/internal/handler"
	"go-todo/internal/repository"
	"go-todo/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/go_todo?sslmode=disable"
	}

	database, err := db.Open(databaseURL)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer database.Close()

	todoRepository := repository.NewPostgresRepository(database)
	todoService := service.NewTodoService(todoRepository)

	fmt.Println("server started at http://localhost:8080")
	if err := handler.NewRouter(todoService).Run(":8080"); err != nil {
		fmt.Println("server error:", err)
	}
}
