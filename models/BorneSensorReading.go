package models

import "time"

type BorneSensorReading struct {
	SensorID    string    `json:"SensorID" format:"uuid" doc:"The ID of the borne sensor."`
	Timestamp   time.Time `json:"timestamp" readOnly:"true" doc:"The moment that this reading was done."`
	Location    Point     `json:"location" doc:"The value read by the location sensor."`
	Heartbeat   *int      `json:"heartbeat,omitempty" doc:"The value read by the heartbeat sensor."`
	Temperature *int      `json:"temperature,omitempty" doc:"The value read by the temperature sensor in degrees Celsius."`
}
