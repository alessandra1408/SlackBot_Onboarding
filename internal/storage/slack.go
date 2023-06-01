package storage

import (
	"github.com/slack-go/slack"
)

type slackStorage struct {
	Api *slack.Client
}

func NewSlackStorage(api *slack.Client) *slackStorage {
	return &slackStorage{
		Api: api,
	}
}

func (s *slackStorage) GetUserByEmail(email string) (*slack.User, error) {
	return s.Api.GetUserByEmail(email)
}

func (s *slackStorage) SendMessage(channel string, messageText string, boolean bool) (string, string, string, error) {
	return s.Api.SendMessage(channel, slack.MsgOptionText(messageText, boolean))
}

func (s *slackStorage) PostMessage(channelID string, messageText string, boolean bool) (string, string, error) {
	return s.Api.PostMessage(channelID, slack.MsgOptionText(messageText, boolean))
}

func (s *slackStorage) DeleteMessage(channel, messageTimestamp string) (string, string, error) {
	return s.Api.DeleteMessage(channel, messageTimestamp)
}

func (s *slackStorage) UpdateMessage(channelID, timestamp string, messageText string, boolean bool) (string, string, string, error) {
	return s.Api.UpdateMessage(channelID, timestamp, slack.MsgOptionText(messageText, boolean))
}
