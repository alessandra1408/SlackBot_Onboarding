package slack

import (
	"log"

	"github.com/slack-go/slack"
)

func NewSlackConnection(slackBotToken string) (*slack.Client, error) {
	api := slack.New(slackBotToken)
	rtm := api.NewRTM()
	_, _, sErr := rtm.ConnectRTM()
	if sErr != nil {
		log.Printf("Error to connect with slack. Err %s\n", sErr)
	}

	return api, sErr
}
