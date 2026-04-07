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
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("error creating rabbitmq connection: %v\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Connection to Server successful!")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("error retrieving username: %v", err)
		return
	}

	_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, "transient")
	if err != nil {
		fmt.Printf("error declaring and binding channel and queue: %v\n", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("\nProgram is shutting down")
}
