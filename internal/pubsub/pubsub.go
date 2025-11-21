package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	SimpleQueueDurable   SimpleQueueType = "durable"
	SimpleQueueTransient SimpleQueueType = "transient"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("unable to marshal struct to json: %v", err)
	}
	if err := ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	}); err != nil {
		return fmt.Errorf("error publishing message: %v", err)
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}
	queue, err := ch.QueueDeclare(queueName, queueType == SimpleQueueDurable, queueType == SimpleQueueTransient, queueType == SimpleQueueTransient, false, nil)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}
	if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}
	return ch, queue, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	chanQueue, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	chanDelivery, err := chanQueue.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range chanDelivery {
			var data T
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}
			handler(data)
			msg.Ack(false)
		}
	}()
	return nil
}
