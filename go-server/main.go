package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Resource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Values      []int  `json:"values"`
}

func seedDB(conn *pgxpool.Pool) error {
	_, err := conn.Query(context.Background(), `CREATE TABLE IF NOT EXISTS resources (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT NOT NULL,
      values INT[] NOT NULL
    );`)

	return err
}

// 2. Criar um canal para receber os dados do json

// 3. criar uma goroutine para receber os dados do canal e inserir no banco
// conn.Query(context.Background(), "insert into resources(name, description, values) values($1, $2, $3)", resource.Name, resource.Description, resource.Values)

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = seedDB(conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to seed database: %v\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	//⭐⭐
	resourceChannel := make(chan Resource)

	//⭐⭐⭐
	go func() {
		for resource := range resourceChannel {
			_, err := conn.Exec(context.Background(), "INSERT INTO resources (name, description, values) VALUES ($1, $2, $3)", resource.Name, resource.Description, resource.Values)
			if err != nil {
				// Lidar com erros de inserção no banco de dados
				fmt.Fprintf(os.Stderr, "Erro na inserção no banco de dados: %v\n", err)
			}
		}
	}()

	r.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		// 1.
		// converter o body para json: https://gowebexamples.com/json/
		// inserir o json no canal para ser inserido no banco
		// retornar código 202 Accepted

		//⭐
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var resource Resource
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&resource); err != nil {
			http.Error(w, "Falha ao decodificar JSON", http.StatusBadRequest)
			return
		}

		resourceChannel <- resource // Insira no canal

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "Accepted")
	}).Methods("POST")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
