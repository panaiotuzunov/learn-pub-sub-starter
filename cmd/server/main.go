package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	fmt.Println("Programs is shutting down...")
	connection.Close()
}
