package models

import (
	"time"
)

type DetectionRecord struct {
	ID         string            `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this detection."`
	SpeciesID  string            `json:"speciesID" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the species of the detected animal."`
	SensorID   string            `json:"sensorID" minLength:"2" doc:"The identification for the sensor that detected the animal."`
	SensorType string            `json:"sensorType" enum:"visual,acoustic,motion,radio,chemical,other" doc:"The type of the sensor that detected the animal."`
	Location   Point             `json:"location" doc:"The location that the animal was detected at."`
	Timestamp  time.Time         `json:"timestamp" format:"date-time" doc:"The moment the animal was detected."`
	URI        *string           `json:"uri,omitempty" format:"uri" doc:"The URI to view this detection in a publicly available external system."`
	Animals    []DetectionAnimal `json:"animals" doc:"Information on the animals involved in this detection."`
}

type Detection struct {
	DetectionRecord
	Species Species `json:"species" doc:"The species of the animal(s) involved in this detection."`
}
