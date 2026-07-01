// handler と DB のやり取りを抽象化するためのインターフェースと、Postgres 用の実装を定義する
// Todo の保存処理に対する共通ルールを定義する
package main

// どのリクエストに紐づいた処理なのか途中でキャンセルされたか、タイムアウトしたかを伝えるための変数
import "context"

// Todo を保存・取得・更新・削除するためには以下の6つのメソッドを持っていてください、というルールを定義する
// 実際のPostgreSQL処理はpostgres_store.goにあります
/* ctx の役割は、リクエストのライフサイクルに応じて処理を中断するための仕組みを提供することです。
例えば、クライアントがリクエストをキャンセルした場合や、タイムアウトが発生した場合に、データベース操作を中断することができます。
これにより、不要なリソース消費を防ぎ、アプリケーションのパフォーマンスと信頼性を向上させることができます。 */
type TodoStore interface {
	ListTodos(ctx context.Context) ([]Todo, error)
	CreateTodo(ctx context.Context, req CreateTodoRequest) (Todo, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	UpdateTodo(ctx context.Context, id int, req UpdateTodoRequest) (Todo, error)
	DeleteTodo(ctx context.Context, id int) error
	ToggleTodoComplete(ctx context.Context, id int) (Todo, error)
}

// TodoStoreのインスタンスを保持する変数を定義する
// service.go では、todoStore.ListTodos()のように、TodoStoreのメソッドを呼び出すことで、DBとのやり取りを行うことができます
var todoStore TodoStore

// TodoServiceのインスタンスを保持する変数を定義する
// handlers.go では、todoService.ListTodos()のように、TodoServiceのメソッドを呼び出すことで、Todoの処理を行うことができます
var todoService *TodoService
