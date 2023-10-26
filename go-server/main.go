package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Resource struct {
	ID          int    `json:"id"`
	Name        string `json:"nome"`
	Description string `json:"descricao"`
	Values      []int  `json:"valores"`
}

// 0.Criar um campo de createdAt do tipo date ou type stamp cirado automaticmaente pelo banco
func seedDB(conn *pgxpool.Pool) error {

	_, err := conn.Query(context.Background(), `DROP TABLE IF EXISTS resources;`)

	if err != nil {
		return err
	}

	_, errCreate := conn.Query(context.Background(), `CREATE TABLE IF NOT EXISTS resources (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT NOT NULL,
      values INT[] NOT NULL
    );`)

	return errCreate
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

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter o argumento para inteiro: %v\n", err)
		os.Exit(1)
	}

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

	if err := decoder.Decode(&dataSet); err == io.EOF {
		fmt.Println("Leitura do Json concluida")
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao decodificar JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Leitura do JSON concluida: %d rows \n", len(dataSet))

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

	fmt.Println("DB semeado com sucesso!")

	var wg sync.WaitGroup
	//4.Criar um for que insere os elementos da lista dentro do resourcechannel
	resourceChannel := make(chan Resource)

	for i := 0; i < nWorkersInt; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for resource := range resourceChannel {
				_, err := conn.Exec(context.Background(), "INSERT INTO resources (name, description, values) VALUES ($1, $2, $3)", resource.Name, resource.Description, resource.Values)
				if err != nil {
					continue
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, resource := range dataSet {
			// Envie o objeto 'data' para o 'resourceChannel'
			resourceChannel <- Resource{
				ID:          resource.ID,
				Name:        resource.Name,
				Description: resource.Description,
				Values:      resource.Values,
			}
		}
		close(resourceChannel)
	}()

	//3.Atraves de um for criar a quantidade de goroutines equivalente
	//a quantidade de workers que executam o canal de dados de maneira simultanea
	wg.Wait()
	os.Exit(0)
}
