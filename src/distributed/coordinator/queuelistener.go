package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/BesWick/PISensor/src/distributed/dto"
	"github.com/BesWick/PISensor/src/distributed/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

//QueueListener discovers the data queues, recieve msgs and
//translate these msgs into events in a EventAggregator
type QueueListener struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	//registry of all sources coordinator is listening on
	sources map[string]<-chan amqp.Delivery //key=queue's Name //val = recieve only channel
}

//NewQueueListener initializes a new QueueListener Obj
func NewQueueListener() *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
	}
	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

//ListenForNewSource allows QueueListener to discover new sensors
func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch)

	//rebinding ch from default exhange to fanout
	ql.ch.QueueBind(q.Name, "", "amq.fanout", false, nil)

	msgs, _ := ql.ch.Consume(q.Name, "", true, false, false, false, nil)

	for msg := range msgs {
		//new sensor coming online
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), //access to specific sensor queue
			"",
			true,
			false,
			false,
			false,
			nil)
		//check if sensor is already registered
		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan
			//add listener to that sensor
			go ql.AddListener(sourceChan)
		}
	}
}

//AddListener adds a queuelistener to read sensor data
func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Recieved message: %v\n", sd)
	}
}
