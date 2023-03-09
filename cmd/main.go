package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shomali11/slacker"
)

type header struct {
	from    string
	to      []string
	subject string
}

type Message struct {
	header      *header
	body        string
	attachments []*file
}

type file struct {
	Name     string
	Header   map[string][]string
	CopyFunc func(w io.Writer) error
}

type FileSetting func(*file)

func main() {
	err := godotenv.Load("/home/alessandra-goncalves/Documents/estudos/Go/SlackBot_Onboarding/.env")
	if err != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	gmailAccount := os.Getenv("GMAIL_ACCOUNT")
	gmailAccountName := os.Getenv("GMAIL_ACCOUNT_NAME")
	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	gmailRecipient := os.Getenv("GMAIL_RECIPIENT")
	gmailRecipientName := os.Getenv("GMAIL_RECIPIENT_NAME")

	bot := slacker.NewClient(botToken, appToken)

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "yob calcutator",
		Examples:    []string{"my yob is 2023"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println("error: ", err)
			}
			age := time.Now().Year() - yob
			r := fmt.Sprintf("age is %d", age)
			sendEmail(gmailAccountName, gmailAccount, gmailPassword, gmailRecipientName, gmailRecipient)
			response.Reply(r)

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Message) appendFile(list []*file, name string, settings []FileSetting) []*file {
	f := &file{
		Name:   filepath.Base(name),
		Header: make(map[string][]string),
		CopyFunc: func(w io.Writer) error {
			h, err := os.Open(name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, h); err != nil {
				h.Close()
				return err
			}
			return h.Close()
		},
	}

	for _, s := range settings {
		s(f)
	}

	if list == nil {
		return []*file{f}
	}

	return append(list, f)
}

func (m *Message) Attach(filename string, settings ...FileSetting) {
	m.attachments = m.appendFile(m.attachments, filename, settings)
}

func sendEmail(gmailAccountName, gmailAccount, gmailPassword, gmailRecipientName, gmailRecipient string) {
	fmt.Println("Gmail Account: ", gmailAccount)
	fmt.Println("Gmail Account Name: ", gmailAccountName)
	fmt.Println("Gmail Password: ", gmailPassword)

	fmt.Println("Gmail Recipient: ", gmailRecipient)
	fmt.Println("Gmail Recipient Name: ", gmailRecipientName)

	// Initialise the required mail message variables
	from := mail.NewEmail(gmailAccountName, gmailAccount)
	subject := "Let's Send an Email With Golang and SendGrid"
	to := mail.NewEmail(gmailRecipientName, gmailRecipient)
	plainTextContent := "Here is your AMAZING email!"
	htmlContent := "Here is your <strong>AMAZING</strong> email!"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Attempt to send the email
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
	}

	// Check if it was sent
	statusCode := response.StatusCode
	if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		fmt.Println("Email sent!")
	}

	if err != nil {
		log.Panicf("Some error occurred with send email. Err %s", err)
	}

}
