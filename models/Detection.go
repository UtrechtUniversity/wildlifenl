package models

import (
	"time"
)

type DetectionRecord struct {
	ID        int       `json:"-"`
	SensorID  string    `json:"sensorID" doc:"The identification for the sensor that detected the animal."`
	Location  Point     `json:"location" doc:"The location that the animal was detected at."`
	Timestamp time.Time `json:"timestamp" format:"date-time" doc:"The moment the animal was detected."`
	SpeciesID string    `json:"speciesID" format:"uuid" writeOnly:"true" required:"true" docs:"The ID of the species of the detected animal."`
}

type Detection struct {
	DetectionRecord
	Species Species `json:"species" doc:"The species of the detected animal."`
}
