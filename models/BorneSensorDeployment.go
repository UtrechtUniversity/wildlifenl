package models

import "time"

type BorneSensorDeploymentRecord struct {
	AnimalID string     `json:"animalID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the animal that bears this sensor."`
	SensorID string     `json:"sensorID" doc:"The ID of the borne sensor."`
	Start    time.Time  `json:"start" format:"date-time" doc:"The moment this deployment started."`
	End      *time.Time `json:"end" format:"date-time" required:"false" readOnly:"true" doc:"The moment this deployment finished."`
}

type BorneSensorDeployment struct {
	BorneSensorDeploymentRecord
	Animal Animal `json:"animal" doc:"The animal that bears this sensor."`
}
