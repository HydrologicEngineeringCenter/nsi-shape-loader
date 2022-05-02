package elevation

import (
	"fmt"
)

type Point struct {
	FdId      int     `db:"fd_id"`
	X         float64 `db:"x"`
	Y         float64 `db:"y"`
	Elevation float64 `default:"-999999"`
}

type Points []*Point

type BoundingBox struct {
	MinX float64 `json:"minX"`
	MaxX float64 `json:"maxX"`
	MinY float64 `json:"minY"`
	MaxY float64 `json:"maxY"`
}

// NilElevation checks whether Point contains elevation data
func (p Point) NilElevation() bool {
	return p.Elevation == -999999
}

// Intersect checks whether a list of Points intersect with a National Map Item
func (p Points) IsIntersecting(i Item) bool {
	for _, point := range p {
		if i.BoundingBox.Contains(*point) {
			return true
		}
	}
	return false
}

// BoundingBox calculates the BoundingBox for a set of Points
func (p Points) BoundingBox() BoundingBox {
	b := BoundingBox{}
	for _, point := range p {
		if point.X > b.MaxX {
			b.MaxX = point.X
		}
		if point.X < b.MinX {
			b.MinX = point.X
		}
		if point.Y > b.MaxY {
			b.MaxY = point.Y
		}
		if point.Y < b.MinY {
			b.MinY = point.Y
		}
	}
	if len(p) == 0 {
		return BoundingBox{}
	} else {
		return b
	}
}

// Contains checks whether Point is within the BoundingBox
func (b BoundingBox) Contains(p Point) bool {
	return b.MinX <= p.X && p.X <= b.MaxX && b.MinY <= p.Y && p.Y <= b.MaxY
}

// Intersect takes a list of Points and filter those not contained within the BoundingBox
func (b BoundingBox) Intersect(p Points) Points {
	var selectedPoints Points
	for _, point := range p {
		if b.Contains(*point) {
			selectedPoints = append(selectedPoints, point)
		}
	}
	return selectedPoints
}

func (b BoundingBox) QueryNationalMap() (QueryResult, error) {
	query := NewNationalMapQuery()
	query.setParam("bbox", fmt.Sprintf("%f,%f,%f,%f", b.MinX, b.MinY, b.MaxX, b.MaxY))
	r, err := query.sendRequest()
	return r, err
}
