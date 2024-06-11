package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Location struct {
	Latitude  float64 `json:"latitude" minimum:"-90" maximum:"90"`
	Longitude float64 `json:"longitude" minimum:"-180" maximum:"180"`
}

func (p Location) Value() (driver.Value, error) {
	return string("(" + strconv.FormatFloat(p.Longitude, 'f', 6, 64) + "," + strconv.FormatFloat(p.Latitude, 'f', 6, 64) + ")"), nil
}

func (p *Location) Scan(value interface{}) error {
	if value == nil { // Questionable situation.
		return nil
	}
	text := string(value.([]uint8))
	text = strings.Trim(text, " \t\n")
	if text[0] == '(' && text[len(text)-1] == ')' {
		text = strings.TrimSuffix(strings.TrimPrefix(text, "("), ")")
		parts := strings.Split(text, ",")
		lon, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 64)
		if err != nil {
			return fmt.Errorf("failed to scan into Point: %w", err)
		}
		lat, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 64)
		if err != nil {
			return fmt.Errorf("failed to scan into Point: %w", err)
		}
		*p = Location{Latitude: lat, Longitude: lon}
		return nil
	}
	return fmt.Errorf("failed to scan into Point")
}
