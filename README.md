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

ファイル分割しました。動作は変えずに、`go-todo2` と同じように `cmd/api` と `internal` 配下へ役割ごとに整理しています。

- [cmd/api/main.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/cmd/api/main.go:1)
サーバー起動、DB 接続、各レイヤーの組み立て

- [internal/model/todo.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/model/todo.go:1)
Todo, request 用 struct

- [internal/handler/todo_handler.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/handler/todo_handler.go:1)
GET / POST / PUT / DELETE / PATCH の handler 本体

- [internal/handler/router.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/handler/router.go:1)
Gin のルーティング登録

- [internal/service/todo_service.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/service/todo_service.go:1)
Todo の入力チェックと処理ルール

- [internal/repository/todo_repository.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/repository/todo_repository.go:1)
PostgreSQL への保存・取得・更新・削除

- [internal/db/db.go (line 1)](/Users/greenpowermarco/personal_files/my_app/go-practice/go-todo/internal/db/db.go:1)
PostgreSQL 接続とテーブル作成

---

## レイヤー分割

現在の実装では、処理を `handler` / `service` / `repository` / `model` / `db` に分けています。

```txt
HTTP request
  ↓
handler
  ↓
service
  ↓
repository
  ↓
PostgreSQL
```

### handler 層

`internal/handler` が handler 層です。

handler は、HTTP リクエストとレスポンスを扱います。

- URL パラメータを取得する
- JSON リクエストを構造体に変換する
- service を呼び出す
- HTTP ステータスコードと JSON レスポンスを返す

Gin を使っているため、handler は `*gin.Context` を受け取ります。

```go
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := getTodoIDFromContext(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := h.service.GetTodo(c.Request.Context(), id)
	// ...
}
```

### service 層

`internal/service` が service 層です。

service は、Todo API の処理ルールを担当します。

- Todo を作成する前に title を確認する
- Todo を更新する前に title を確認する
- repository を呼び出して保存・取得・更新・削除する

例えば、`title` が空かどうかの確認は handler ではなく service 側で行います。

```go
func (s *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error) {
	if req.Title == "" {
		return model.Todo{}, ErrTitleRequired
	}

	return s.repository.CreateTodo(ctx, req)
}
```

### バリデーション

現在の実装では、`internal/service/todo_service.go` で Todo 作成・更新時の入力値をチェックしています。

- `title` は前後の空白を取り除いたうえで必須
- `title` は100文字以内
- `description` は前後の空白を取り除く
- `description` は500文字以内

handler は JSON の読み取りとレスポンスを担当し、入力値のルールは service に集めています。

```go
func validateCreateTodoRequest(req model.CreateTodoRequest) (model.CreateTodoRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		return model.CreateTodoRequest{}, ErrTitleRequired
	}

	if len([]rune(req.Title)) > MaxTitleLength {
		return model.CreateTodoRequest{}, ErrTitleTooLong
	}

	return req, nil
}
```

### repository 層

`internal/repository/todo_repository.go` が repository 層です。

repository は、データの保存先とのやり取りを担当します。

- PostgreSQL から Todo 一覧を取得する
- PostgreSQL に Todo を登録する
- PostgreSQL の Todo を更新・削除する

保存処理に必要なメソッドを `TodoRepository` interface として定義しています。

```go
type TodoRepository interface {
	ListTodos(ctx context.Context) ([]model.Todo, error)
	CreateTodo(ctx context.Context, req model.CreateTodoRequest) (model.Todo, error)
	GetTodo(ctx context.Context, id int) (model.Todo, error)
	UpdateTodo(ctx context.Context, id int, req model.UpdateTodoRequest) (model.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
	ToggleTodoComplete(ctx context.Context, id int) (model.Todo, error)
}
```

実際の PostgreSQL 処理は `PostgresRepository` にあります。

### レイヤー分割のメリット

- handler が HTTP の処理に集中できる
- service に処理ルールを集められる
- repository を差し替えやすくなる
- テストで PostgreSQL ではなくメモリ実装を使いやすい
- コードの責務が分かりやすくなる

今回のテストでは、PostgreSQL ではなく `memory_store_test.go` のメモリ実装を使っています。handler は service を呼ぶだけなので、本番は PostgreSQL、テストはメモリ、という差し替えがしやすくなっています。

---

## テストを書く

Go では、`*_test.go` という名前のファイルにテストを書きます。

実行コマンドは以下です。

```bash
go test ./...
```

`./...` は、現在のディレクトリ以下のすべてのパッケージをテストするという意味です。

### このプロジェクトのテスト構成

現在は、主に2種類のテストを書いています。

```txt
main_test.go
  router + handler のテスト

internal/service/todo_service_test.go
  service 層の単体テスト
```

### handler テスト

`main_test.go` では、HTTP リクエストに近い形で handler をテストしています。

```go
req := httptest.NewRequest(http.MethodPost, "/todos", body)
rec := httptest.NewRecorder()

newTestRouter(store).ServeHTTP(rec, req)
```

このテストでは、以下を確認できます。

- ルーティングが正しく動くか
- HTTP ステータスコードが正しいか
- JSON レスポンスが正しいか
- エラー時のレスポンス形式が正しいか

Gin を使っている場合、handler 関数を直接呼ぶより、router 経由でリクエストする方が自然です。

### service テスト

`internal/service/todo_service_test.go` では、HTTP を通さずに service 層だけをテストしています。

このテストでは、以下を確認できます。

- `title` や `description` のバリデーションが正しいか
- 前後の空白が取り除かれるか
- バリデーションエラー時に repository が呼ばれないか
- repository から返る `ErrTodoNotFound` がそのまま扱えるか
- 完了状態の切り替え処理が正しいか

service テストでは、PostgreSQL ではなく fake repository を使っています。

```go
repo := newFakeTodoRepository(nil, 1)
service := NewTodoService(repo)
```

これにより、DB を起動しなくても service の処理ルールだけを高速に確認できます。

### テストの種類

このプロジェクトでは、次のように考えると整理しやすいです。

- unit test: service のように、1つの層だけを小さく確認する
- handler test: HTTP リクエストに近い形で、router と handler を確認する
- integration test: 実際の PostgreSQL に接続して、repository まで含めて確認する

現在追加しているのは、unit test と handler test です。

PostgreSQL を使った integration test は、Docker やテスト用DBの準備が必要になるため、次の段階で追加するとよいです。

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
go run ./cmd/api
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
r.GET("/todos", todoHandler.GetTodos)
r.POST("/todos", todoHandler.CreateTodo)
r.GET("/todos/:id", todoHandler.GetTodo)
r.PUT("/todos/:id", todoHandler.UpdateTodo)
r.DELETE("/todos/:id", todoHandler.DeleteTodo)
r.PATCH("/todos/:id/complete", todoHandler.CompleteTodo)
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
var req model.CreateTodoRequest
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
