package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/go_todo?sslmode=disable"
	}

	db, err := openDB(databaseURL)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	// PostgreSQLとやり取りするStoreを作成し、TodoStore interfaceとして保持する
	todoStore = NewPostgresStore(db)
	// StoreをServiceに渡し、handlerからはService経由でTodoの処理を呼び出す
	todoService = NewTodoService(todoStore)

	fmt.Println("server started at http://localhost:8080")
	if err := newRouter().Run(":8080"); err != nil {
		fmt.Println("server error:", err)
	}
}
