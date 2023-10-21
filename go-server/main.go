package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Resource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Values      []int  `json:"values"`
}

func seedDB(conn *pgx.Conn) (pgx.Rows, error) {
	rows, err := conn.Query(context.Background(), `CREATE TABLE IF NOT EXISTS resources (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT NOT NULL,
      values INT[] NOT NULL
    );`)

	return rows, err
}

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	_, err = seedDB(conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to seed database: %v\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	r.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}).Methods("POST")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
