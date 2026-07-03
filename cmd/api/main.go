package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	authService := service.NewAuthService(
		getEnv("AUTH_USERNAME", "admin"),
		getEnv("AUTH_PASSWORD", "password"),
		getEnv("JWT_SECRET", "go-todo-dev-secret"),
	)
	allowedOrigins := getCSVEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	fmt.Println("server started at http://localhost:8080")
	if err := handler.NewRouter(todoService, authService, allowedOrigins).Run(":8080"); err != nil {
		fmt.Println("server error:", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getCSVEnv(key, fallback string) []string {
	value := getEnv(key, fallback)
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
