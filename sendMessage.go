package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

const (
	envVarPath = "/home/alessandra-goncalves/Documents/neoway/key-results/SlackBot/.env"
)

var (
	// jiraEmail         = os.Getenv("JIRA_MAIL")
	// jiraUrl           = os.Getenv("JIRA_INSTANCE")
	// jiraToken         = os.Getenv("JIRA_AUTH_TOKEN")
	slackBotTokenVar = "SLACK_BOT_TOKEN"
	slackAppTokenVar = "SLACK_APP_TOKEN"
	// projectKeyDefault = os.Getenv("JIRA_PROJECT_KEY_DEFAULT")
)

type Person struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func main() {

	err := sendMessageToUser()
	if err != nil {
		log.Fatalf("Error sending message to user. Err: %v", err)
	}

}

func getSlackAuth() (botToken, appToken string, bot *slacker.Slacker, api *slack.Client, err error) {
	err = godotenv.Load(envVarPath)
	if err != nil {
		return "", "", nil, nil, err
	}

	botToken = os.Getenv(slackBotTokenVar)
	if botToken == "" {
		return "", "", nil, nil, fmt.Errorf("%s environment variable is not set", slackBotTokenVar)
	}
	appToken = os.Getenv(slackAppTokenVar)
	if appToken == "" {
		return "", "", nil, nil, fmt.Errorf("%s environment variable is not set", slackAppTokenVar)
	}

	bot = slacker.NewClient(botToken, appToken)
	api = slack.New(botToken)

	/* bot.BotCommands(){

	} */

	return botToken, appToken, bot, api, nil
}

// sendMessageToUser function sends an onboarding message to the user
func sendMessageToUser() error {
	_, _, bot, api, err := getSlackAuth()
	if err != nil {
		return err
	}

	bot.Command("Mensagem para <email>", &slacker.CommandDefinition{
		Description: "Send message of onboarding to new coworkers",
		Examples:    []string{"Fazer onboarding da @aleh"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			personEmail := request.StringParam("email", "null")
			personEmail = strings.ReplaceAll(personEmail, "<mailto:", "")
			personEmail = strings.ReplaceAll(personEmail, ">", "")
			formatedEmail := strings.Split(personEmail, "|")
			personName, personID, err := getUserInfo(api, string(formatedEmail[0]))
			if err != nil {
				log.Printf("Some error occurred in getUserInfo function. Err %s\n", err)
				return
			}

			mensagemOnboarding := fmt.Sprintf(`Olá %s, agora você faz parte do squad Sebrae :slightly_smiling_face:.
Para o seu processo de onboarding, temos diversos materiais de gestão de conhecimento e apoio no confluence. É muito importante que você entre nesse espaço (também é onde armazenamos as nossas documentações).
Já irei te passar alguns links úteis para o seu dia a dia (alguns desses você provavelmente não tem acesso, mas eles já estão sendo solicitados):
			1. %v (%v)
			2. %v (%v)
			3. %v (%v e %v)
			4. %v (%v)
			5. %v (%v)`, personName, "Confluence SEBRAE", "https://tinyurl.com/confluencesebrae", "Bucket", "https://tinyurl.com/bucketSEB", "GCP de QA e PROD", "https://tinyurl.com/GCPqaSEB", "https://tinyurl.com/GCPsebPROD", "Projeto de Service Ops (para acessar o datalake SEBRAE)", "https://tinyurl.com/ProjServOps", "Épico central do Squad no Jira", "https://tinyurl.com/jiraSEB")

			_, _, err = api.PostMessage(
				personID,
				slack.MsgOptionText(mensagemOnboarding, false),
			)
			if err != nil {
				log.Printf("Some error occurred in postMessage to user. Err %s", err)
				return
			}

			str := fmt.Sprintf("Mensagem de onboarding enviada para %v!", personName)
			err = response.Reply(str)
			if err != nil {
				log.Printf("Some error occurred in sendMessageToUser Function. Err %s", err)
				return
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	if err != nil {
		log.Fatalf("Some error occurred in defer cancel(). Err %v", err)
	}
	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// getUserInfo function returns the person's name and ID based on the email address
func getUserInfo(api *slack.Client, personEmail string) (personName, personID string, err error) {
	err = godotenv.Load(envVarPath)
	if err != nil {
		return "", "", err
	}
	user, err := api.GetUserByEmail(personEmail)
	if err != nil {
		log.Printf("Some error occurred in getUserInfo internal function. Err %s\n", err)
	}

	return user.Profile.RealName, user.ID, err
}
