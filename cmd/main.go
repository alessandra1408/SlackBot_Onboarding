package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alessandra1408/SlackBot_Onboarding/domain/core"
	"github.com/alessandra1408/SlackBot_Onboarding/internal/storage"
	database "github.com/alessandra1408/SlackBot_Onboarding/internal/storage/dataBase"
	slack "github.com/alessandra1408/SlackBot_Onboarding/internal/storage/slack"
	"github.com/alessandra1408/SlackBot_Onboarding/pkg/env"
	"github.com/joho/godotenv"
)

const (
	slackBotTokenVar        = "SLACK_BOT_TOKEN"
	slackAppTokenVar        = "SLACK_APP_TOKEN"
	projectKeyDefault       = "JIRA_PROJECT_KEY_DEFAULT"
	jiraToken               = "JIRA_AUTH_TOKEN"
	jiraInstance            = "JIRA_INSTANCE"
	postgresHost            = "ENV_POSTGRES_HOST"
	postgresPort            = "ENV_POSTGRES_PORT"
	postgresUser            = "ENV_POSTGRES_USER"
	postgresScrt            = "ENV_POSTGRES_SCRT"
	postgresDataBase        = "ENV_POSTGRES_DATABASE"
	postgresSchema          = "ENV_POSTGRES_SCHEMA"
	sslMode                 = "SSL_MODE"
	envPostgresTimezone     = "ENV_POSTGRES_TIMEZONE"
	defaultPostgresTimezone = "UTC"
)

//init inicializa as variaveis de ambiente
func init() {
	if err := godotenv.Load(filepath.Join("/home/alessandra-goncalves/Documents/neoway/key-results/SlackBot/", ".env")); err != nil {
		log.Print("Error to load enviroment variables")
	}
}

func main() {
	a, b := getSlackBotToken(), getSlackAppToken()
	fmt.Println(a, b)

	env.CheckRequired(slackAppTokenVar, slackBotTokenVar, projectKeyDefault, jiraToken, jiraInstance, postgresHost, postgresPort, postgresUser, postgresScrt, postgresDataBase, postgresSchema)

	//criando conexão com o banco
	db, dbErr := database.NewConnection(getPostgresURL())
	if dbErr != nil {
		log.Fatalf("Error on open connection with postgres. Err %s", dbErr)
		os.Exit(1)
	}

	//instanciando o banco
	pg := storage.NewPostgresStorage(db)

	// pg.Queryx("SELECT id, name, email, last_mr_date FROM list_of_people")

	//criando conexão com o slack e client
	slackClient, sErr := slack.NewSlackConnection(getSlackAppToken())
	if sErr != nil {
		log.Fatalf("Error on open connection with slack. Err %s", sErr)
		os.Exit(1)
	}

	//instanciando o slack
	slack := storage.NewSlackStorage(slackClient)

	service := core.NewService(pg, slack)

	service.Running()

}

func getSlackBotToken() string {
	return env.GetString(slackBotTokenVar)
}

func getSlackAppToken() string {
	return env.GetString(slackAppTokenVar)
}

func getJiraToken() string {
	return env.GetString(jiraToken)
}

func getJiraIntance() string {
	return env.GetString(jiraInstance)
}

func getProjectKeyDefault() string {
	return env.GetString(projectKeyDefault)
}

func getPostgresHost() string {
	return env.GetString(postgresHost)
}

func getPostgresPort() string {
	return env.GetString(postgresPort)
}

func getPostgresUser() string {
	return env.GetString(postgresUser)
}

func getPostgresScrt() string {
	return env.GetString(postgresScrt)
}

func getPostgresDatabase() string {
	return env.GetString(postgresDataBase)
}

func getPostgresSchema() string {
	return env.GetString(postgresSchema)
}

func getPostgresSSLMode() string {
	return env.GetString(sslMode)
}

func getPostgresTimezone() string {
	return env.GetString(envPostgresTimezone, defaultPostgresTimezone)
}

func getPostgresURL() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s timezone=%s",
		getPostgresHost(),
		getPostgresPort(),
		getPostgresDatabase(),
		getPostgresUser(),
		getPostgresScrt(),
		getPostgresSSLMode(),
		getPostgresTimezone())
}
