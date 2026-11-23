package models

import "time"

type BorneSensorDeploymentRecord struct {
	ID                     string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this borne-sensor deployment."`
	AnimalID               string     `json:"animalID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the animal that bears this sensor."`
	SensorID               string     `json:"sensorID" minLength:"2" doc:"The ID of the borne sensor."`
	ContactHardwareAddress *string    `json:"contactHardwareAddress" pattern:"^([0-9A-F]{2}:){5}[0-9A-F]{2}$" doc:"The EUI-48 hardware address of the bluetooth contact tracing device."`
	Start                  time.Time  `json:"start" format:"date-time" doc:"The moment this deployment started."`
	End                    *time.Time `json:"end" format:"date-time" required:"false" readOnly:"true" doc:"The moment this deployment finished."`
}

type BorneSensorDeployment struct {
	BorneSensorDeploymentRecord
	Animal              Animal               `json:"animal" doc:"The animal that bears this sensor."`
	User                User                 `json:"user" doc:"The user that deployed this sensor."`
	BorneSensorReadings []BorneSensorReading `json:"borneSensorReadings,omitempty" readOnly:"true" doc:"The borneSensorReadings associated with this deployment."`
}
