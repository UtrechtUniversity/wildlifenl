package wildlifenl

import "time"

type Operations struct {
	Endpoint string
}

type SpatiotemporalInput struct {
	Start     time.Time `query:"start" format:"date-time" required:"true" doc:"The start moment of the spatiotemporal span."`
	End       time.Time `query:"end" format:"date-time" required:"true" doc:"The end moment of the spatiotemporal span."`
	Latitude  float64   `query:"latitude" minimum:"-90" maximum:"90" required:"true" doc:"The latitude of the spatiotemporal centroid."`
	Longitude float64   `query:"longitude" minimum:"-180" maximum:"180" required:"true" doc:"The longitude of the spatiotemporal centroid."`
	Radius    int       `query:"radius" minimum:"1" maximum:"10000" required:"true" doc:"The approximate radius of the spatiotemporal span in meters relative to the centroid."`
}
