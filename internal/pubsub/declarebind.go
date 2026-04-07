package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	Durable   SimpleQueueType = "durable"
	Transient SimpleQueueType = "transient"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("error creating channel form connection: %v\n", err)
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	var newQueue amqp.Queue
	if queueType == "durable" {
		newQueue, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	} else if queueType == "transient" {
		newQueue, err = ch.QueueDeclare(queueName, false, true, true, false, nil)
	} else {
		fmt.Printf("not a valid queue type\n")
		return &amqp.Channel{}, amqp.Queue{}, nil
	}
	if err != nil {
		fmt.Printf("error creaing new queue: %v", err)
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	err = ch.QueueBind(newQueue.Name, key, exchange, false, nil)
	if err != nil {
		fmt.Print("error binding new queue to channel: %v", err)
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	return ch, newQueue, nil
}
