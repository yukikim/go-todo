# GO Practice
## Todo REST API

作る機能:
Todo 一覧取得
GET /todos

Todo 詳細取得
GET /todos/{id}

Todo 作成
POST /todos

Todo 更新
PUT /todos/{id}

Todo 削除
DELETE /todos/{id}

完了状態の切り替え
PATCH /todos/{id}/complete

**DB を使わず、メモリ上のスライスや map で管理する**

```go
var todos = map[int]Todo{}
```

## 基本的なサーバーを作成
**標準ライブラリだけで Todo API を作る**

### 1.プロジェクト作成
```bash
mkdir go-todo
cd go-todo
go mod init go-todo
```
### 2.main.go を作成
```bash
touch main.go
```
### 3.最初はこの構成で始める
```txt
go-todo/
  go.mod
  main.go
```

### 4.main.go に最低限のサーバーを書く
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	fmt.Println("server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
```

### 5.サーバー起動
```bash
go run main.go
```

### 6.別ターミナルで確認
```bash
curl http://localhost:8080/health
```
**ok が返れば成功**

---
## Todo 構造体を作る

```go
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
```

## メモリ上に map[int]Todo を用意する

```go
var todos = map[int]Todo{}
var nextID = 1
```