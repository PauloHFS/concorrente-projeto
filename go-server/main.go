package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Resource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Values      []int  `json:"values"`
}

//0.Criar um campo de createdAt do tipo date ou type stamp cirado automaticmaente pelo banco
func seedDB(conn *pgxpool.Pool) error {
	_, err := conn.Query(context.Background(), `CREATE TABLE IF NOT EXISTS resources (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT NOT NULL,
      values INT[] NOT NULL
    );`)

	return err
}


func main() {
//1.Pegar a quantidade de workers como argumento
	nWorkers := os.Args[1]
	if _, err := strconv.Atoi(arg); err == nil {
		fmt.Println("O argumento é um número inteiro.")
	} else {
		fmt.Println("O argumento não é um número inteiro.")
		os.exit(1)
	}
//2.Ler os dados do arquivo do dataset e colocar em uma lista
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

//3.Atraves de um for criar a quantidade de goroutines equivalente
//a quantidade de workers que executam o canal de dados de maneira simultanea
	go func() {
		for resource := range resourceChannel {
			_, err := conn.Exec(context.Background(), "INSERT INTO resources (name, description, values) VALUES ($1, $2, $3)", resource.Name, resource.Description, resource.Values)
			if err != nil {
				// Lidar com erros de inserção no banco de dados
				fmt.Fprintf(os.Stderr, "Erro na inserção no banco de dados: %v\n", err)
			}
		}
	}()

//4.Criar um for que insere os elementos da lista dentro do resourcechannel
resourceChannel := make(chan Resource)

}
