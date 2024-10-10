package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Polygon []Point

func (p Polygon) Value() (driver.Value, error) {
	values := make([]string, 0)
	for _, point := range p {
		value, err := point.Value()
		if err != nil {
			return nil, err
		}
		values = append(values, value.(string))
	}
	return string("(" + strings.Join(values, ",") + ")"), nil
}

func (p *Polygon) Scan(value interface{}) error {
	if value == nil { // Questionable situation.
		return nil
	}
	polygon := make([]Point, 0)
	text := strings.ReplaceAll(string(value.([]uint8)), " ", "")
	if strings.HasPrefix(text, "((") && strings.HasSuffix(text, "))") {
		text = strings.TrimSuffix(strings.TrimPrefix(text, "(("), "))")
		pointTexts := strings.Split(text, "),(")
		for _, pointText := range pointTexts {
			point := new(Point)
			if err := point.Scan([]uint8("(" + pointText + ")")); err != nil {
				return err
			}
			polygon = append(polygon, *point)
		}
		*p = polygon
		return nil
	}
	return fmt.Errorf("failed to scan into Polygon")
}
