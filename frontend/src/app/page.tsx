"use client";

import { FormEvent, useCallback, useMemo, useState } from "react";
import styles from "./page.module.css";

type Todo = {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  createdAt: string;
  updatedAt: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export default function Home() {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("password");
  const [token, setToken] = useState("");
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const completedCount = useMemo(
    () => todos.filter((todo) => todo.completed).length,
    [todos],
  );

  const fetchTodos = useCallback(async (currentToken: string) => {
    setError("");

    try {
      const response = await fetch(`${API_URL}/todos`, {
        headers: {
          Authorization: `Bearer ${currentToken}`,
        },
      });

      if (!response.ok) {
        throw new Error("Todoを取得できませんでした");
      }

      const data: Todo[] = await response.json();
      setTodos(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    }
  }, []);

  async function login(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        throw new Error("ログインできませんでした");
      }

      const data: { token: string } = await response.json();
      setToken(data.token);
      await fetchTodos(data.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  }

  async function createTodo(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(`${API_URL}/todos`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, description }),
      });

      if (!response.ok) {
        throw new Error("Todoを作成できませんでした");
      }

      setTitle("");
      setDescription("");
      await fetchTodos(token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  }

  async function toggleTodo(id: number) {
    setError("");

    try {
      const response = await fetch(`${API_URL}/todos/${id}/complete`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("完了状態を更新できませんでした");
      }

      await fetchTodos(token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    }
  }

  async function deleteTodo(id: number) {
    setError("");

    try {
      const response = await fetch(`${API_URL}/todos/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Todoを削除できませんでした");
      }

      await fetchTodos(token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    }
  }

  function logout() {
    setToken("");
    setTodos([]);
    setError("");
  }

  return (
    <main className={styles.page}>
      <section className={styles.header}>
        <div>
          <p className={styles.kicker}>Go Todo</p>
          <h1>Todo Dashboard</h1>
        </div>
        {token ? (
          <button className={styles.secondaryButton} onClick={logout}>
            ログアウト
          </button>
        ) : null}
      </section>

      {error ? <p className={styles.error}>{error}</p> : null}

      {!token ? (
        <form className={styles.panel} onSubmit={login}>
          <label>
            ユーザー名
            <input
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              autoComplete="username"
            />
          </label>
          <label>
            パスワード
            <input
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              type="password"
              autoComplete="current-password"
            />
          </label>
          <button className={styles.primaryButton} disabled={isLoading}>
            {isLoading ? "ログイン中..." : "ログイン"}
          </button>
        </form>
      ) : (
        <>
          <section className={styles.stats}>
            <div>
              <span>{todos.length}</span>
              <p>すべて</p>
            </div>
            <div>
              <span>{completedCount}</span>
              <p>完了</p>
            </div>
            <div>
              <span>{todos.length - completedCount}</span>
              <p>未完了</p>
            </div>
          </section>

          <form className={styles.panel} onSubmit={createTodo}>
            <label>
              タイトル
              <input
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                placeholder="例: Next.jsからAPIを呼ぶ"
              />
            </label>
            <label>
              説明
              <textarea
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                placeholder="内容を入力"
                rows={3}
              />
            </label>
            <button className={styles.primaryButton} disabled={isLoading}>
              {isLoading ? "保存中..." : "追加"}
            </button>
          </form>

          <section className={styles.list}>
            {todos.map((todo) => (
              <article className={styles.todo} key={todo.id}>
                <div>
                  <p className={todo.completed ? styles.doneTitle : ""}>
                    {todo.title}
                  </p>
                  {todo.description ? <span>{todo.description}</span> : null}
                </div>
                <div className={styles.actions}>
                  <button
                    className={styles.secondaryButton}
                    onClick={() => toggleTodo(todo.id)}
                  >
                    {todo.completed ? "戻す" : "完了"}
                  </button>
                  <button
                    className={styles.dangerButton}
                    onClick={() => deleteTodo(todo.id)}
                  >
                    削除
                  </button>
                </div>
              </article>
            ))}
          </section>
        </>
      )}
    </main>
  );
}
