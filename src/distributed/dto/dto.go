package dto

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name      string
	Value     float64
	TimeStamp time.Time
}

//using gob for encoding msg - good for pure go app
//JSON's a better choice - more language neutral
func init() {
	gob.Register(SensorMessage{})
}
