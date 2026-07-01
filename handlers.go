package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := todoStore.ListTodos(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list todos")
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストのJSONを格納するための構造体を作成(CreateTodoRequestはTypescriptのinterfaceのようなもの)
	var req CreateTodoRequest
	// json.NewDecoder(r.Body)でリクエストボディをデコードし、reqに格納する(&reqの&はreqの住所を渡すことで、関数内でreqの値を変更できるようにするため)
	// Goでは、関数に引数として渡すときに値をコピーするため、関数内で変更しても元の変数には影響しません。&を使うことで、変数のアドレスを渡し、関数内で変更した値が元の変数に反映されるようにしています
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// JSONのデコードに失敗した場合、HTTPステータスコード400(Bad Request)を返す
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// タイトルが空の場合Todoとして登録しない、HTTPステータスコード400(Bad Request)を返す
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := todoStore.CreateTodo(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create todo")
		return
	}

	// 作成したTodoをJSON形式でレスポンスとして返すために、writeJSON関数を呼び出しています
	writeJSON(w, http.StatusCreated, todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := todoStore.GetTodo(r.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get todo")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func getTodoIDFromRequest(r *http.Request) (int, error) {
	// chi.URLParam(r, "id")でURLパラメータからidを取得し、strconv.Atoiで文字列を整数に変換しています
	idText := chi.URLParam(r, "id")
	return strconv.Atoi(idText)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo ID")
		return
	}

	_, err = todoStore.GetTodo(r.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get todo")
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := todoStore.UpdateTodo(r.Context(), id, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update todo")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo ID")
		return
	}

	err = todoStore.DeleteTodo(r.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := todoStore.ToggleTodoComplete(r.Context(), id)
	if errors.Is(err, errTodoNotFound) {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update todo")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}
