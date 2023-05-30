package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const envVarPath = "/home/alessandra-goncalves/Documents/neoway/key-results/SlackBot/.env"

type personAvailable struct {
	ID           int
	Name         string
	Email        string
	Last_MR_Date string
}

func main() {
	err := godotenv.Load(envVarPath)
	if err != nil {
		fmt.Printf("Some error occured in read .env file. Err %s", err)
	}

	postgresHost := os.Getenv("ENV_GCP_POSTGRES_HOST")
	postgresPort := os.Getenv("ENV_GCP_POSTGRES_PORT")
	postgresUser := os.Getenv("ENV_GCP_POSTGRES_USER")
	postgresScrt := os.Getenv("ENV_GCP_POSTGRES_SCRT")
	postgresDatabase := os.Getenv("ENV_GCP_POSTGRES_DATABASE")
	// postgresSchema := os.Getenv("ENV_GCP_POSTGRES_SCHEMA")
	postgresURL := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", postgresHost,
		postgresPort, postgresDatabase, postgresUser, postgresScrt)

	db, err := NewConnection(postgresURL)
	if err != nil {
		log.Printf("error on PostgreSQL connection in main function: %q", err)
	}

	// close database
	defer db.Close()
	if err != nil {
		log.Printf("error on Close PostgreSQL connection: %q", err)
	}

	// db.Query("insert into list_of_people (name,	email, last_mr_date) values ('Rudson Souza ', 'rudson.souza@neoway.com.br', '10/05/2023');")

	// Executar a consulta SQL
	rows, err := db.Query("SELECT id, name, email, last_mr_date FROM list_of_people")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Processar os resultados da consulta
	var people []personAvailable
	for rows.Next() {
		var person personAvailable
		err := rows.Scan(&person.ID, &person.Name, &person.Email, &person.Last_MR_Date)
		if err != nil {
			panic(err.Error())
		}
		people = append(people, person)
	}

	// Usar o resultado da consulta em outras partes do código
	for _, person := range people {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Last_MR_Date: %s\n", person.ID, person.Name, person.Email, person.Last_MR_Date)
	}

}

func NewConnection(postgresURL string) (*sql.DB, error) {

	db, sErr := sql.Open("postgres", postgresURL)
	if sErr != nil {
		log.Printf("error on PostgreSQL connection: %q", sErr)
		return nil, sErr
	}
	pErr := db.Ping()
	if pErr != nil {
		log.Printf("error on Ping DB: %q", pErr)
		return nil, pErr
	}
	return db, nil
}

func CheckCriticalError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckError(content string, err error) {
	if err != nil {
		fmt.Printf("Some error occurred with %v. Err %v\n", content, err)
	}
}
