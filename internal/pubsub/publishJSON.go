package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		fmt.Printf("error marshaling data: %v\n", err)
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		fmt.Printf("error publishing channel: %v\n", err)
		return err
	}

	return nil
}
