package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

const (
	envVarPath = "/home/alessandra-goncalves/Documents/neoway/key-results/SlackBot/.env"
)

type Project struct {
	Key string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

type CustomField struct {
	Id string `json:"id"`
}

type TextContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Content struct {
	Content []TextContent `json:"content"`
	Type    string        `json:"type"`
}

type Description struct {
	Content []Content `json:"content"`
	Type    string    `json:"type"`
	Version int       `json:"version"`
}

type Fields struct {
	Summary          string      `json:"summary"`
	Issuetype        IssueType   `json:"issuetype"`
	Description      Description `json:"description"`
	Project          Project     `json:"project"`
	Customfield10298 CustomField `json:"customfield_10298"`
	Customfield12100 CustomField `json:"customfield_12100"`
}

type Task struct {
	Fields Fields `json:"fields"`
}

func main() {
	err := createJiraIssue()
	if err != nil {
		log.Fatalf("Some error occurred in createJiraIssue. Err %v", err)
	}
}

func getSlackTokens() (string, string, error) {
	err := godotenv.Load(envVarPath)
	if err != nil {
		log.Printf("error reading .env file. Err %v\n", err)
		return "", "", err
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackAppToken := os.Getenv("SLACK_APP_TOKEN")

	if slackBotToken == "" {
		return "", "", fmt.Errorf("%s environment variable is not set", slackBotToken)
	}
	if slackAppToken == "" {
		return "", "", fmt.Errorf("%s environment variable is not set", slackAppToken)
	}

	return slackBotToken, slackAppToken, err
}

func getSlackAuth(slackBotToken, slackAppToken string) (*slacker.Slacker, *slack.Client) {

	bot := slacker.NewClient(slackBotToken, slackAppToken)
	api := slack.New(slackBotToken)

	return bot, api
}

func getJiraConfig() (string, string, string, string, error) {
	err := godotenv.Load(envVarPath)
	if err != nil {
		log.Printf("error reading .env file. Err %v\n", err)
		return "", "", "", "", err
	}

	jiraUrl := os.Getenv("JIRA_INSTANCE")
	jiraEmail := os.Getenv("JIRA_MAIL")
	jiraToken := os.Getenv("JIRA_AUTH_TOKEN")
	projectKeyDefault := os.Getenv("JIRA_PROJECT_KEY_DEFAULT")

	return jiraUrl, jiraEmail, jiraToken, projectKeyDefault, err
}

func setNewTask(summary, description, projectKey string) *Task {
	return &Task{
		Fields: Fields{
			Summary:   summary,
			Issuetype: IssueType{Name: "Consulting Service"},
			Description: Description{
				Content: []Content{
					{Content: []TextContent{
						{Text: description, Type: "text"},
					},
						Type: "paragraph",
					},
				},
				Type:    "doc",
				Version: 1,
			},
			Project: Project{
				Key: projectKey,
			},
			Customfield10298: CustomField{
				Id: "17316",
			},
			Customfield12100: CustomField{
				Id: "18816",
			},
		},
	}
}

func createJiraIssue() error {
	slackBotToken, SlackAppToken, err := getSlackTokens()
	if err != nil {
		fmt.Printf("erro em ler getSlackAuth. Err %v\n", err)
		return err
	}

	bot, _ := getSlackAuth(slackBotToken, SlackAppToken)

	jiraUrl, jiraEmail, jiraToken, projectKeyDefault, err := getJiraConfig()
	if err != nil {
		fmt.Printf("erro em ler getJiraConfig. Err %v\n", err)
		return err
	}

	bot.Command("task <context1> <context2>", &slacker.CommandDefinition{
		Description: "Send message of onboarding to new coworkers",
		Examples:    []string{"Fazer onboarding da @aleh"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			summary := request.Param("context1")
			description := request.StringParam("context2", "Descrição Não informada pelo usuário")
			// summary := request.StringParam("summary", "null")
			// Description := request.StringParam("Description", "null")

			summary = strings.Trim(summary, "\"")
			description = strings.Trim(description, "\"")

			if summary == "" {
				response.Reply("Os Parametros obrigatórios não foram passados corretamente. Task não criada.")
				log.Fatalf("Os Parametros não foram passados corretamente. Task não criada.")
			}

			newIssue := setNewTask(summary, description, projectKeyDefault)

			issueJson, err := json.Marshal(newIssue)
			if err != nil {
				log.Printf("some error occurred. Err %v", err)
			}

			fmt.Println("issueJson: ", string(issueJson))

			req, err := http.NewRequest("POST", jiraUrl, bytes.NewBuffer(issueJson))
			if err != nil {
				log.Printf("some error occurred. Err %v", err)
			}

			req.SetBasicAuth(jiraEmail, jiraToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			respBody, err := io.ReadAll(resp.Body)

			var body []map[string]string
			json.Unmarshal(respBody, &body)
			// fmt.Println("var body[5]: ", body[5])

			fmt.Println("resp.Status", resp.StatusCode)
			fmt.Println("resp.Body", string(respBody))
			defer resp.Body.Close()
			if err != nil {
				log.Fatalf("Some error occured in close resp.Body Function. Err %s", err)
			}
			str := "Task criada!"
			err = response.Reply(str)
			if err != nil {
				log.Printf("Some error occured in createJiraIssue Function. Err %s", err)
				return
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
