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

	gameState := gamelogic.NewGameState(userName)

	ch, err := connection.Channel()
	if err != nil {
		fmt.Printf("error creating channel: %v\n", err)
		return
	}

	pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, "pause."+userName, routing.PauseKey, "transient", handlerPause(gameState))

	pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+userName, routing.ArmyMovesPrefix+".*", "transient", handlerMove(gameState, ch))

	pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+".*", "durable", handlerWar(gameState, ch))

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		if input[0] == "spawn" {
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("error spawning unit: %v\n", err)
			}
		} else if input[0] == "move" {
			armyMove, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("error moving unit: %v\n", err)
			}
			fmt.Printf("Move of %v to %v was successful\n", armyMove.Units, armyMove.ToLocation)
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, "army_moves."+userName, armyMove)
			if err != nil {
				fmt.Printf("error publishing move to JSON: %v\n", err)

			} else {
				fmt.Printf("Move Published Successfully\n")
			}
		} else if input[0] == "status" {
			gameState.CommandStatus()
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Printf("Spamming not allowed yet!\n")
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Printf("Command Not Understood...\nTry Again...\n")
			continue
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("\nProgram is shutting down")
}
