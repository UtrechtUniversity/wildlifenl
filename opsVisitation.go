package wildlifenl

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

type VisitationHolder struct {
	Body *models.Visitation `json:"visitation"`
}

type VisitationInput struct {
	Start       time.Time `query:"start" format:"date-time" required:"true" doc:"The start moment of the time span."`
	End         time.Time `query:"end" format:"date-time" required:"true" doc:"The end moment of the time span."`
	LivingLabID string    `query:"livingLabID" format:"uuid" required:"true" doc:"The ID of the living lab."`
	CellSize    float64   `query:"cellSize" minimum:"20" maximum:"10000" required:"true" doc:"The size of the heatmap cells in meters."`
}

type visitationOperations Operations

func newVisitationOperations() *visitationOperations {
	return &visitationOperations{Endpoint: "visitation"}
}

func (o *visitationOperations) RegisterGetForLivingLab(api huma.API) {
	name := "Get Visitation For Living Lab"
	description := "Retrieve visitation within a time span for a certain living lab."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"nature-area-manager", "wildlife-manager", "herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *VisitationInput) (*VisitationHolder, error) {
		livingLab, err := stores.NewLivingLabStore(relationalDB).Get(input.LivingLabID)
		if err != nil {
			return nil, handleError(err)
		}
		livingLabPoints := make([]orb.Point, 0)
		for _, p := range livingLab.Definition {
			livingLabPoints = append(livingLabPoints, orb.Point{p.Longitude, p.Latitude})
		}
		livingLabPoints = append(livingLabPoints, orb.Point{livingLab.Definition[0].Longitude, livingLab.Definition[0].Latitude})

		trackingReadings, err := stores.NewTrackingReadingStore(relationalDB, timeseriesDB).GetForTimespan(input.Start, input.End)
		if err != nil {
			return nil, handleError(err)
		}
		humanPoints := make([]orb.Point, 0)
		for _, r := range trackingReadings {
			humanPoints = append(humanPoints, orb.Point{r.Location.Longitude, r.Location.Latitude})
		}

		mercatorPolygon := project.Polygon(orb.Polygon{livingLabPoints}, project.WGS84.ToMercator)
		b := mercatorPolygon.Bound()
		minX := b.Min[0]
		minY := b.Min[1]
		maxX := b.Max[0]
		maxY := b.Max[1]

		nx := int(math.Ceil((maxX - minX) / input.CellSize))
		ny := int(math.Ceil((maxY - minY) / input.CellSize))
		type cell struct {
			centroidX float64
			centroidY float64
			count     int
		}
		cells := make(map[int]cell)
		key := func(x, y int) int { return x*1000000 + y }

		for ix := range nx {
			centroidX := minX + (float64(ix)+0.5)*input.CellSize
			for iy := range ny {
				centroidY := minY + (float64(iy)+0.5)*input.CellSize
				if planar.RingContains(mercatorPolygon[0], orb.Point{centroidX, centroidY}) { // Keep cells whose centroid is inside the polygon.
					cells[key(ix, iy)] = cell{centroidX: centroidX, centroidY: centroidY, count: 0}
				}
			}
		}

		for _, p := range humanPoints { // Aggregate human points into cells.
			mercatorPoint := project.Point(p, project.WGS84.ToMercator)
			x := mercatorPoint[0]
			y := mercatorPoint[1]
			ix := int(math.Floor((x - minX) / input.CellSize))
			iy := int(math.Floor((y - minY) / input.CellSize))
			if ix < 0 || iy < 0 || ix >= nx || iy >= ny {
				continue
			}
			k := key(ix, iy)
			if gc, ok := cells[k]; ok {
				gc.count++
				cells[k] = gc
			}
		}

		visitation := models.Visitation{}
		for _, c := range cells {
			centroidWGS := project.Point(orb.Point{c.centroidX, c.centroidY}, project.Mercator.ToWGS84)
			visitation = append(visitation, models.VisitationCell{
				Centroid: models.Point{Latitude: centroidWGS[1], Longitude: centroidWGS[0]},
				Count:    c.count,
			})
		}
		return &VisitationHolder{Body: &visitation}, nil
	})
}
