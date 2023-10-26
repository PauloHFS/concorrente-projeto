package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Resource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Values      []int  `json:"values"`
}

// 0.Criar um campo de createdAt do tipo date ou type stamp cirado automaticmaente pelo banco
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
	if _, err := strconv.Atoi(nWorkers); err == nil {
		fmt.Println("O argumento é um número inteiro.")
	} else {
		fmt.Println("O argumento não é um número inteiro.")
		os.Exit(1)
	}
	nWorkersInt, err := strconv.Atoi(nWorkers)

	//2.Ler os dados do arquivo do dataset e colocar em uma lista
	file, err := os.Open("dataset.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao abrir o arquivo: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	//2.1 Crie um decodificador JSON
	decoder := json.NewDecoder(file)

	//2.2 Declare uma fatia para armazenar os elementos do JSON
	var dataSet []Resource

	for {
		var resource Resource
		if err := decoder.Decode(&resource); err == io.EOF {
			break // Chegou ao final do arquivo
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao decodificar JSON: %v\n", err)
			os.Exit(1)
		}

		// Adicione o objeto 'resource' à fatia 'dataSet'
		dataSet = append(dataSet, resource)
	}

	//0. área que conecta o database
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

	//4.Criar um for que insere os elementos da lista dentro do resourcechannel
	resourceChannel := make(chan Resource)

	for _, resource := range dataSet {
		// Envie o objeto 'data' para o 'resourceChannel'
		resourceChannel <- Resource{
			ID:          resource.ID,
			Name:        resource.Name,
			Description: resource.Description,
			Values:      resource.Values,
		}
	}

	//3.Atraves de um for criar a quantidade de goroutines equivalente
	//a quantidade de workers que executam o canal de dados de maneira simultanea
	for i := 0; i < nWorkersInt; i++ {
		go func() {
			for resource := range resourceChannel {
				_, err := conn.Exec(context.Background(), "INSERT INTO resources (name, description, values) VALUES ($1, $2, $3)", resource.Name, resource.Description, resource.Values)
				if err != nil {
					// Lidar com erros de inserção no banco de dados
					fmt.Fprintf(os.Stderr, "Erro na inserção no banco de dados: %v\n", err)
				}
			}
		}()
	}

}
