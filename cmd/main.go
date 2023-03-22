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

type Person struct {
	ID       string
	Name     string
	Email    string
	Password string
}

var env = "/home/alessandra-goncalves/Documents/estudos/Go/SlackBot_Onboarding/.env"

func main() {
	err := sendMessageToUser()
	if err != nil {
		log.Fatalf("Error sending message to user. Err: %v", err)
	}
}

// auth function returns the required variables for the program
func auth() (botToken, appToken string, bot *slacker.Slacker, api *slack.Client, err error) {
	err = godotenv.Load(env)
	fmt.Println(err)
	if err != nil {
		return "", "", nil, nil, err
	}

	botToken = os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return "", "", nil, nil, fmt.Errorf("SLACK_BOT_TOKEN environment variable is not set")
	}
	appToken = os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return "", "", nil, nil, fmt.Errorf("SLACK_APP_TOKEN environment variable is not set")
	}

	bot = slacker.NewClient(botToken, appToken)
	api = slack.New(botToken)

	return botToken, appToken, bot, api, nil
}

// config function returns the person's password from the .env file
func config() Person {
	err := godotenv.Load(env)
	if err != nil {
		log.Fatalf("Error loading .env file. Err: %s\n", err)
	}

	var person Person
	person.Password = os.Getenv("GMAIL_PASSWORD")

	return person
}

// sendMessageToUser function sends an onboarding message to the user
func sendMessageToUser() error {
	_, _, bot, api, err := auth()
	if err != nil {
		return err
	}

	bot.Command("Mensagem <email>", &slacker.CommandDefinition{
		Description: "Send message of onboarding to new coworkers",
		Examples:    []string{"Fazer onboarding da @aleh"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			personEmail := request.StringParam("email", "null")
			personEmail = strings.ReplaceAll(personEmail, "<mailto:", "")
			personEmail = strings.ReplaceAll(personEmail, ">", "")
			formatedEmail := strings.Split(personEmail, "|")
			fmt.Println("personEmail: ", formatedEmail[0])
			personName, personID, err := getUserInfo(api, formatedEmail[0])
			if err != nil {
				log.Printf("Some error occured in getUserInfo function. Err %s\n", err)
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
				log.Printf("Some error occured in postMessage to user. Err %s", err)
				return
			}

			str := fmt.Sprintf("Mensagem de onboarding enviada para %v!", personName)
			err = response.Reply(str)
			if err != nil {
				log.Printf("Some error occured in sendMessageToUser Function. Err %s", err)
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

// getUserInfo function returns the person's name and ID based on the email address
func getUserInfo(api *slack.Client, personEmail string) (personName, personID string, err error) {
	err = godotenv.Load(env)
	if err != nil {
		return "", "", err
	}
	user, err := api.GetUserByEmail(personEmail)
	if err != nil {
		log.Printf("Some error occured in GetUserInfo function. Err %s\n", err)
	}

	return user.Profile.RealName, user.ID, err
}
