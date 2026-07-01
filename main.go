package main

import (
	"fmt"
	"log"
	"net/http"
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

	// TodoStoreのインスタンスをPostgresStoreで初期化する(./postgres_store.goで定義)
	todoStore = NewPostgresStore(db)

	fmt.Println("server started at http://localhost:8080")
	// server を起動してエラーがなければ、newRouter()(./router.go) でルーティングを設定したハンドラを渡す
	if err := http.ListenAndServe(":8080", newRouter()); err != nil {
		fmt.Println("server error:", err)
	}
}
