package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func main() {
	err := godotenv.Load("/home/alessandra-goncalves/Documents/estudos/Go/SlackBot_CalculeteAge/.env")
	if err != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(botToken, appToken)

	go printCommandEvents(bot.CommandEvents())

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

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println("Timestamp: ", event.Timestamp)
		fmt.Println("Command: ", event.Command)
		fmt.Println("Parameters: ", event.Parameters)
		fmt.Println("Event: ", event.Event)
	}
}
