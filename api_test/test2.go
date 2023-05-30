package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

const (
	envVarPath = "/home/alessandra-goncalves/Documents/neoway/key-results/SlackBot/.env"
)

func main() {

	err := godotenv.Load(envVarPath)
	if err != nil {
		log.Fatal("Error to read .env file")
	}
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// Defina o canal e o modal (payload) que deseja enviar
	channelID := os.Getenv("CHANNEL_ID")
	payloadJSON := `{
		"blocks": [
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Olá, este é um exemplo de modal em formato JSON!"
				}
			}
		]
	}`

	_, _, err = api.PostMessage(channelID, slack.MsgOptionText(" ", false), slack.MsgOptionAttachments(
		slack.Attachment{
			Text: payloadJSON,
		},
	))

	if err != nil {
		log.Fatalf("Erro ao enviar a mensagem: %v", err)
	}
}
