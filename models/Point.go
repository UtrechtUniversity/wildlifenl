package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	Latitude  float64 `json:"latitude" minimum:"-90" maximum:"90"`
	Longitude float64 `json:"longitude" minimum:"-180" maximum:"180"`
}

func (p Point) Value() (driver.Value, error) {
	return string("(" + strconv.FormatFloat(p.Longitude, 'f', 6, 64) + "," + strconv.FormatFloat(p.Latitude, 'f', 6, 64) + ")"), nil
}

func (p *Point) Scan(value interface{}) error {
	if value == nil { // Questionable situation.
		return nil
	}
	text := strings.ReplaceAll(string(value.([]uint8)), " ", "")
	if text[0] == '(' && text[len(text)-1] == ')' {
		text = strings.TrimSuffix(strings.TrimPrefix(text, "("), ")")
		parts := strings.Split(text, ",")
		longitude, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return fmt.Errorf("failed to scan into Point: %w", err)
		}
		latitude, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return fmt.Errorf("failed to scan into Point: %w", err)
		}
		*p = Point{Latitude: latitude, Longitude: longitude}
		return nil
	}
	return fmt.Errorf("failed to scan into Point")
}
