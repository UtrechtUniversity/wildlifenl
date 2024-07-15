package models

import "time"

type BorneSensorDeploymentRecord struct {
	AnimalID string     `json:"AnimalID,omitempty" format:"uuid" writeOnly:"true" doc:"The ID of the animal that bears this sensor."`
	SensorID string     `json:"SensorID" format:"uuid" doc:"The ID of the borne sensor."`
	Start    time.Time  `json:"start" doc:"The moment this deployment started."`
	End      *time.Time `json:"end" required:"false" readOnly:"true" doc:"The moment this deployment finished."`
}

type BorneSensorDeployment struct {
	BorneSensorDeploymentRecord
	Animal Animal `json:"animal" doc:"The animal that bears this sensor."`
}
