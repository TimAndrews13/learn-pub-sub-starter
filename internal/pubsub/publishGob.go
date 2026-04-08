package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		fmt.Printf("error encoding val: %v\n", err)
		return err
	}

	data := buf.Bytes()

	err := ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        data,
		})
	if err != nil {
		fmt.Printf("error publishing channel: %v\n", err)
		return err
	}

	return nil
}
