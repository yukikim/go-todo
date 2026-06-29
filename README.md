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

**以下の実装状況はコミットを参照**

---

#### Todo 登録(POST)
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Goを学習する","description":"POST /todosを作る"}'
```
#### Todo 取得(GET)
```bash
curl http://localhost:8080/todos/1
```
#### Todo 更新(PUT)
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"新しいタイトル","description":"新しい説明","completed":true}'
```
#### Todo 削除(DELETE)
```bash
curl -X DELETE http://localhost:8080/todos/1
```
#### 完了状態を切り替える処理
```bash
curl -X PATCH http://localhost:8080/todos/1/complete
```

---

ファイル分割しました。動作は変えずに、役割ごとに整理しています。

- [main.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/main.go:1)
サーバー起動とルーティング登録だけ

- [models.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/models.go:1)
Todo, request 用 struct, メモリ保存用の todos / nextID

- [handlers.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/handlers.go:1)
GET / POST / PUT / DELETE / PATCH の handler 本体

- [response.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/response.go:1)
writeJSON, writeError, ErrorResponse

#### 起動コマンド
```bash
go run .
```

#### 確認コマンド
```bash
curl http://localhost:8080/health
curl http://localhost:8080/todos
```