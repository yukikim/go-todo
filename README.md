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
Todo, request 用 struct

- [handlers.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/handlers.go:1)
GET / POST / PUT / DELETE / PATCH の handler 本体

- [response.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/response.go:1)
writeJSON, writeError, ErrorResponse

---

## PostgreSQL に保存する

現在の実装では、Todo をメモリ上の map ではなく PostgreSQL に保存します。

### PostgreSQL 起動

```bash
docker compose up -d postgres
```

### 環境変数

未指定の場合は、以下の接続先を使います。

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_todo?sslmode=disable
```

別の接続先を使う場合は、起動前に `DATABASE_URL` を指定してください。

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/go_todo?sslmode=disable"
```

#### 起動コマンド

```bash
go run .
```

#### 確認コマンド
```bash
curl http://localhost:8080/health
curl http://localhost:8080/todos
```

---

## Gin を使ったルーティング

現在の実装では、標準ライブラリの `http.HandleFunc` ではなく、`Gin` を使ってルーティングしています。

Gin は Go 製の Web フレームワークです。標準ライブラリだけでも API は作れますが、Gin を使うとルーティング、パスパラメータ取得、JSON の入力・出力などを短く書けます。

### 導入パッケージ

```bash
go get github.com/gin-gonic/gin
```

### ルーティング例

```go
r.GET("/todos", getTodosHandler)
r.POST("/todos", createTodoHandler)
r.GET("/todos/:id", getTodoHandler)
r.PUT("/todos/:id", updateTodoHandler)
r.DELETE("/todos/:id", deleteTodoHandler)
r.PATCH("/todos/:id/complete", completeTodoHandler)
```

標準ライブラリだけで書く場合は、`/todos` と `/todos/` を分けて登録したり、handler の中で `switch r.Method` を使って HTTP メソッドを判定したりする必要がありました。

Gin では `GET`, `POST`, `PUT`, `DELETE`, `PATCH` をルーティング定義で分けられるため、handler の責務がシンプルになります。

### メリット

- ルーティングを読みやすく書ける
- `/todos/:id` のようにパスパラメータを扱いやすい
- JSON リクエストを `ShouldBindJSON` で構造体に変換できる
- JSON レスポンスを `c.JSON` で返せる
- 存在しないルートや許可していないメソッドの処理をまとめて定義できる
- middleware を追加しやすい

### 今回の実装で変わった点

ID の取得は、文字列操作ではなく Gin の `Param` を使います。

```go
idText := c.Param("id")
```

JSON リクエストの読み取りは、`json.NewDecoder` ではなく `ShouldBindJSON` を使います。

```go
var req CreateTodoRequest
if err := c.ShouldBindJSON(&req); err != nil {
	writeError(c, http.StatusBadRequest, "invalid JSON")
	return
}
```

JSON レスポンスは `c.JSON` で返します。

```go
c.JSON(http.StatusCreated, todo)
```

### 注意点

- 標準ライブラリだけの実装より依存パッケージが増える
- Gin 独自の `*gin.Context` に依存するため、handler の書き方が標準の `http.HandlerFunc` とは変わる
- テストでは handler を直接呼ぶより、Gin の router 経由でリクエストする方が自然
- 小さな API では標準ライブラリだけでも十分な場合がある
- フレームワークの便利さに頼りすぎると、HTTP の基本処理が見えにくくなることがある

学習目的では、まず標準ライブラリで API の仕組みを理解し、その後に Gin を導入すると違いが分かりやすいです。
