package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func InsertUser(db *sql.DB) error {
	names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"}

	for _, name := range names {
		var exists bool
		queryCheck := `SELECT EXISTS (SELECT 1 FROM users WHERE name = $1)`
		err := db.QueryRow(queryCheck, name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking user %s: %w", name, err)
		}

		if !exists {
			queryInsert := `INSERT INTO users (name) VALUES ($1)`
			_, err := db.Exec(queryInsert, name)
			if err != nil {
				return fmt.Errorf("error inserting user %s: %w", name, err)
			}
		}
	}
	return nil
}

func ListUsers(db *sql.DB) ([]User, error) {
	if err := CreateUserTable(db); err != nil {
		return nil, err
	}

	if err := InsertUser(db); err != nil {
		return nil, err
	}

	query := `SELECT id, name FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

var startedAt = time.Now()

func main() {
	connStr := "postgres://myuser:mypassword@postgresql-h.api-test-namespace-dev.svc.cluster.local:5432/api-test?sslmode=disable&connect_timeout=10"
	maxAttempts := 5
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Tentativa %d: erro ao abrir a conexão: %v", attempt, err)
		}

		if err == nil {
			// Se a conexão foi bem-sucedida, interrompe o loop
			fmt.Println("Conexão estabelecida com sucesso!")
			break
		}

		log.Printf("Tentativa %d: erro ao conectar ao PostgreSQL: %v", attempt, err)
	}

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		users, err := ListUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/healthz", Healthz)
	http.HandleFunc("/", Hello)
	err = http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
	w.Write([]byte("Hello"))
	w.WriteHeader(http.StatusOK)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(startedAt)

	if duration.Seconds() < 12 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
