package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/BesWick/PISensor/src/distributed/dto"
	"github.com/BesWick/PISensor/src/distributed/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5, "maximum value for generated readings")
var min = flag.Float64("min", 1, "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)

	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)

	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"",               //exchange string,
		sensorQueue.Name, //key string,
		false,            //mandatory bool,
		false,            //immediate bool,
		msg)              //msg amqp.Publishing)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			TimeStamp: time.Now(),
		}

		buff.Reset()
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buff.Bytes(),
		}

		ch.Publish(
			"",             //exchange string,
			dataQueue.Name, //key string,
			false,          //mandatory bool,
			false,          //immediate bool,
			msg)            //msg amqp.Publishing)

		log.Printf("Reading send, Value: %v", value)
	}
}

func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
