package models

type Visitation struct {
	LivingLab *LivingLab       `json:"livingLab" doc:"The living lab associated with this visitation map."`
	Cells     []VisitationCell `json:"cells" doc:"The cells of the visitation map."`
}

type VisitationCell struct {
	Centroid Point `json:"centroid" doc:"The central point of the cell."`
	Count    int   `json:"count" minimum:"0" doc:"The number of human location tracking records in this cell."`
}
