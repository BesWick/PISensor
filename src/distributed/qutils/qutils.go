package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"

//GetChannel gets references to Channel and Connection
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failonError(err, "failed to connect to rabbitMQ")
	ch, err := conn.Channel()
	failonError(err, "Failed to open a channel")

	return conn, ch
}

//GetQueue declares a queue
func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(name, false, false, false, false, nil)
	failonError(err, "Failed to declare a queue")
	return &q
}

func failonError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
