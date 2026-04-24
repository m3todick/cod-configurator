package main

import (
        "database/sql"
        "fmt"
        "log"
        "net/http"
        "os"

        _ "github.com/lib/pq"
)

func main() {
        loadEnv(".env")

        dbURL := os.Getenv("SUPABASE_DB_URL")
        if dbURL == "" {
                log.Fatal("SUPABASE_DB_URL environment variable is required")
        }

        db, err := sql.Open("postgres", dbURL)
        if err != nil {
                log.Fatalf("Ошибка подключения к БД: %v", err)
        }
        defer db.Close()

        if err := db.Ping(); err != nil {
                log.Fatalf("Ошибка подключения к БД: %v", err)
        }
        log.Println("Подключение к БД успешно")

        store := NewStore(db)
        store.clearAllSessions()
        log.Println("Все сохранённые сессии очищены — потребуется повторный вход")

        srv := NewServer(store)

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        fmt.Printf("ЦОД Конфигуратор запущен на http://localhost:%s\n", port)
        if err := http.ListenAndServe(":"+port, srv); err != nil {
                log.Fatal(err)
        }
}
