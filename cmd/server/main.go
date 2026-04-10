package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	gamelogic.PrintServerHelp()

	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("error creating rabbitmq connection: %v\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Connection to Server successful!")

	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("error creating channel from connection: %v\n", err)
		return
	}

	err = pubsub.SubscribeGob(connection, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", "durable", handlerLogs())
	if err != nil {
		fmt.Printf("error subscribing to game log channel: %v", err)
		return
	}

	err = pubsub.PublishJSON(channel, routing.ExchangePerilTopic, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("error publishing json: %v\n", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		if input[0] == "pause" {
			fmt.Printf("Pausing Peril Game...\n")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Printf("error publishing json: %v\n", err)
			}
		} else if input[0] == "resume" {
			fmt.Printf("Resuming Peril Game...\n")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Printf("error publishing json: %v\n", err)
			}
		} else if input[0] == "quit" {
			fmt.Printf("Exiting Peril Game\n")
			break
		} else {
			fmt.Printf("Command Not Understood...\nTry Again...\n")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("\nProgram is shutting down")
}
