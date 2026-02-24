package models

type Visitation []VisitationCell

type VisitationCell struct {
	Centroid Point `json:"centroid"`
	Count    int   `json:"count"`
}
