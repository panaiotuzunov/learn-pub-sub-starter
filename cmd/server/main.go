package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Connection to RabbitMQ established")

	connChan, err := connection.Channel()
	if err != nil {
		log.Fatalf("error creating AMQP channel: %v", err)
	}
	defer connChan.Close()
	if err := pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	fmt.Println("Program is shutting down...")
}
