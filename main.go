package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var todos = map[int]Todo{}
var nextID = 1

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoHandler)

	fmt.Println("server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodosHandler(w, r)
	case http.MethodPost:
		createTodoHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodoHandler(w, r)
	case http.MethodPut:
		updateTodoHandler(w, r)
	case http.MethodDelete:
		deleteTodoHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Todoのスライスを作成し、todosマップの長さを初期容量として設定 make(型、　長さ、　容量)
	todoList := make([]Todo, 0, len(todos)) // todosはTodoのマップ(連想配列)なので、len(todos)で要素数を取得できる
	// rangeは、配列やスライス、マップなどの要素を順番に取り出すための構文です
	for _, todo := range todos { // _, は「無視する」という意味(idは使わない)で、todoはtodosマップの中身だけを順番に取り出すための変数です
		// todosマップの各要素を順番に取り出し、todoListスライスに追加しています
		// マップのまま返すと、JSONのキーが文字列になるため、スライスに変換して返すことで、JSONの配列として返すことができます
		todoList = append(todoList, todo)
	}

	// todoListをJSON形式でレスポンスとして返すために、writeJSON関数を呼び出しています
	writeJSON(w, http.StatusOK, todoList)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストのJSONを格納するための構造体を作成(CreateTodoRequestはTypescriptのinterfaceのようなもの)
	var req CreateTodoRequest
	// json.NewDecoder(r.Body)でリクエストボディをデコードし、reqに格納する(&reqの&はreqの住所を渡すことで、関数内でreqの値を変更できるようにするため)
	// Goでは、関数に引数として渡すときに値をコピーするため、関数内で変更しても元の変数には影響しません。&を使うことで、変数のアドレスを渡し、関数内で変更した値が元の変数に反映されるようにしています
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// JSONのデコードに失敗した場合、HTTPステータスコード400(Bad Request)を返す
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// タイトルが空の場合Todoとして登録しない、HTTPステータスコード400(Bad Request)を返す
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	now := time.Now() // 現在時刻を取得
	// Todo構造体を作成し、ID、タイトル、説明、完了状態、作成日時、更新日時を設定
	todo := Todo{
		ID:          nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 作成したTodoをメモリ上のtodosマップに追加し、次のIDをインクリメント
	todos[todo.ID] = todo
	nextID++

	// 作成したTodoをJSON形式でレスポンスとして返すために、writeJSON関数を呼び出しています
	writeJSON(w, http.StatusCreated, todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, ok := todos[id]
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func getTodoIDFromPath(path string) (int, error) {
	idText := strings.TrimPrefix(path, "/todos/")
	return strconv.Atoi(idText)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, ok := todos[id]
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.Completed = req.Completed
	todo.UpdatedAt = time.Now()
	todos[id] = todo

	writeJSON(w, http.StatusOK, todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	if _, ok := todos[id]; !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	delete(todos, id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
