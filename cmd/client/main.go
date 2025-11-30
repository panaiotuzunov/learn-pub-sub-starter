package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connStr = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection to RabbitMQ established")
	connChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("error creating AMQP channel: %v", err)
	}
	defer connChan.Close()
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not get username: %v", err)
	}
	gs := gamelogic.NewGameState(username)
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs)); err != nil {
		log.Fatalf("Could not subscribe to pause messages: %v", err)
	}
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(gs, connChan)); err != nil {
		log.Fatalf("Could not subscribe to move messages: %v", err)
	}
	if err := pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		"war",
		routing.WarRecognitionsPrefix+".*",
		pubsub.SimpleQueueDurable,
		handlerWar(gs, connChan)); err != nil {
		log.Fatalf("Could not subscribe to war messages: %v", err)
	}
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			if err := gs.CommandSpawn(input); err != nil {
				fmt.Printf("Could not spawn units: %v", err)
			}
		case "move":
			move, err := gs.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := pubsub.PublishJSON(connChan,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+move.Player.Username,
				move); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Move published successfully.")
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) < 2 {
				fmt.Println("Spam command accepts an additional argument.")
				continue
			}
			n, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Println("Invalid number.")
				continue
			}
			for range n {
				log := gamelogic.GetMaliciousLog()
				if err := pubsub.PublishGob(connChan,
					routing.ExchangePerilTopic,
					routing.GameLogSlug+"."+username,
					routing.GameLog{
						CurrentTime: time.Now(),
						Message:     log,
						Username:    username,
					}); err != nil {
					fmt.Printf("error publishing message: %v", err)
				}
			}
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unrecognized command. Please try again.")
			continue
		}
	}
}
