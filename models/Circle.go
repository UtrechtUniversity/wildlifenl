package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Circle struct {
	Location Point   `json:"location"`
	Radius   float64 `json:"radius" minimum:"1" doc:"The radius in meters."`
}

func (c Circle) Value() (driver.Value, error) {
	v, err := c.Location.Value()
	if err != nil {
		return nil, err
	}
	radius := c.Radius * 0.00001 // Converting meters to degrees using 1m = 0.00001 degrees.
	return "<" + v.(string) + "," + strconv.FormatFloat(radius, 'f', 6, 64) + ">", nil
}

func (c *Circle) Scan(value interface{}) error {
	if value == nil { // Questionable situation.
		return nil
	}
	text := strings.ReplaceAll(string(value.([]uint8)), " ", "")
	if text[0] == '<' && text[len(text)-1] == '>' {
		text = strings.TrimSuffix(strings.TrimPrefix(text, "<"), ">")
		parts := strings.Split(text, ")")
		point := new(Point)
		if err := point.Scan([]uint8(parts[0] + ")")); err != nil {
			return fmt.Errorf("failed to scan into Circle[Point]: %w", err)
		}
		radius, err := strconv.ParseFloat(parts[1][1:], 64)
		if err != nil {
			return fmt.Errorf("failed to scan into Circle[Radius]: %w", err)
		}
		radius = radius * 100000 // Converting degrees to meters using 1m = 0.00001 degrees.
		*c = Circle{Location: *point, Radius: radius}
		return nil
	}
	return fmt.Errorf("failed to scan into Circle")
}
