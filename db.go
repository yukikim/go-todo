package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func openDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL) // データベースに接続するためのsql.DBオブジェクトを作成
	if err != nil {                              // データベース接続に失敗した場合はエラーを返す
		return nil, err
	}

	if err := db.Ping(); err != nil { // データベース接続が有効かどうかを確認するためにPingメソッドを呼び出す
		db.Close() // 接続が無効な場合はデータベース接続を閉じる
		return nil, err
	}

	if err := createTodosTable(db); err != nil { // todosテーブルが存在しない場合は作成する
		db.Close() // テーブル作成に失敗した場合はデータベース接続を閉じる
		return nil, err
	}

	return db, nil
}

func createTodosTable(db *sql.DB) error {
	// todosテーブルが存在しない場合は作成するためのSQL文を実行
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			completed BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}
