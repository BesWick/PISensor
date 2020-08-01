package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"
const SensorDiscoveryExchange = "SensorDiscovery" //name of fanout exchange for sensors discovery

//GetChannel gets references to Channel and Connection from RabbitMQ URL
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failonError(err, "failed to connect to rabbitMQ")
	ch, err := conn.Channel()
	failonError(err, "Failed to open a channel")

	return conn, ch
}

//GetQueue returns reference to a new queue in a channel
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
