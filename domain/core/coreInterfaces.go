package core

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/slack-go/slack"
)

type SSlack interface {
	GetUserByEmail(email string) (*slack.User, error)
	SendMessage(channel string, messageText string, boolean bool) (string, string, string, error)
	PostMessage(channelID string, messageText string, boolean bool) (string, string, error)
	DeleteMessage(channel, messageTimestamp string) (string, string, error)
	UpdateMessage(channelID, timestamp string, messageText string, boolean bool) (string, string, string, error)
}

type SService interface {
	GetUserByEmail(email string) (*slack.User, error)
	SendMessage(channel string, messageText string, boolean bool) (string, string, string, error)
	PostMessage(channelID string, messageText string, boolean bool) (string, string, error)
	DeleteMessage(channel, messageTimestamp string) (string, string, error)
	UpdateMessage(channelID, timestamp string, messageText string, boolean bool) (string, string, string, error)
	Queryx(query string) (*sqlx.Rows, error)
	Exec(query string) (sql.Result, error)
	MustExec(query string, args ...any) sql.Result
	Close() error
}

type SPostgresStorage interface {
	Queryx(query string) (*sqlx.Rows, error)
	Exec(query string) (sql.Result, error)
	MustExec(query string, args ...any) sql.Result
	Close() error
}
