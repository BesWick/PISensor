package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	server()
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}

	ch.Publish("", q.Name, false, false, msg)
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {

	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failonError(err, "failed to connect to rabbitMQ")
	ch, err := conn.Channel()
	failonError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failonError(err, "Failed to declare a queue")
	return conn, ch, &q
}

func failonError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
