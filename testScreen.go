package main

import (
	"log"
	"os"

	_ "encoding/json"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go handleIncomingMessage(api, ev)
		}
	}
}

func handleIncomingMessage(api *slack.Client, event *slack.MessageEvent) {
	if event.Text == "exibir_tela" {
		go displayInteractiveScreen(api, event)
	}
}

func displayInteractiveScreen(api *slack.Client, event *slack.MessageEvent) {
	attachment := slack.Attachment{
		Text: "Olá! Clique no botão abaixo.",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "meu_botao",
				Text:  "Clique Aqui",
				Type:  "button",
				Value: "button_click",
			},
		},
	}

	params := slack.MsgOptionAttachments(attachment)

	_, _, err := api.PostMessage(event.Channel, params)
	if err != nil {
		log.Println("Erro ao exibir a tela interativa:", err)
	}
}
