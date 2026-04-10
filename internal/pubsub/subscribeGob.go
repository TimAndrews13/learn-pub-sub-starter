package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
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
		defer ch.Close()
		for d := range deliveries {
			var msg T
			buf := bytes.NewBuffer(d.Body)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(&msg); err != nil {
				fmt.Printf("error decoding Gob: %v", err)
				continue
			}
			ackType := handler(msg)
			switch ackType {
			case Ack:
				d.Ack(false)
				fmt.Println("DEBUG: Ack")
			case NackRequeue:
				d.Nack(false, true)
				fmt.Println("DEBUG: NackRequeue")
			case NackDiscard:
				d.Nack(false, false)
				fmt.Println("DEBUG: NackDiscard")
			}
		}
	}()

	return nil
}
