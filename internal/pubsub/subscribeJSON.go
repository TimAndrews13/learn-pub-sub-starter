package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T)) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		fmt.Printf("error declaring channel and binding queue: %v\n", err)
		return err
	}

	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("error returning deliveries: %v\n", err)
		return err
	}

	go func() {
		for d := range deliveries {
			var msg T
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("error unmarsahlling JSON: %v", err)
				return
			}
			handler(msg)
			err = d.Ack(false)
			if err != nil {
				fmt.Printf("error removing message from queue: %v", err)
				return
			}
		}
		err = ch.Close()
		if err != nil {
			fmt.Printf("error closing channel: %v", err)
			return
		}
	}()

	return nil
}
