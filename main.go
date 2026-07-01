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

	// NewPostgresStore関数(./postgres_store.go)を呼び出して、PostgresStoreのインスタンスを作成し、todoStore変数に代入する
	todoStore = NewPostgresStore(db)

	fmt.Println("server started at http://localhost:8080")
	if err := newRouter().Run(":8080"); err != nil {
		fmt.Println("server error:", err)
	}
}
