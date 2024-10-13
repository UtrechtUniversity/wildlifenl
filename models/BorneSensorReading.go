package models

import "time"

type BorneSensorReadingRecord struct {
	SensorID    string    `json:"sensorID" doc:"The ID of the borne sensor."`
	Timestamp   time.Time `json:"timestamp" format:"date-time" doc:"The moment that this reading was done."`
	Location    *Point    `json:"location,omitempty" doc:"The value read by the location sensor."`
	Altitude    *float64  `json:"altitude,omitempty" doc:"The value read by the altitude sensor."`
	Temperature *float64  `json:"temperature,omitempty" doc:"The value read by the temperature sensor in degrees Celsius."`
	Accelero    *Accelero `json:"accelero,omitempty" doc:"The values read by the accelerometer."`
}

type BorneSensorReading struct {
	BorneSensorReadingRecord
	UserID string `json:"userID" doc:"The ID of the data system that provided this reading."`
}
