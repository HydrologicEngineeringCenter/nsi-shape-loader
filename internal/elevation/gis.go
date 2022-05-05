package elevation

import (
	"errors"

	"github.com/lukeroth/gdal"
)

type Point struct {
	FdId      int      `db:"fd_id"`
	X         float64  `db:"x"`
	Y         float64  `db:"y"`
	Elevation *float64 `db:"ground_elev"` // pointer instead of value for nullable type
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
	return p.Elevation == nil
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
	if len(p) == 0 {
		return BoundingBox{}
	}
	for i, point := range p {
		if i == 0 {
			b.MinX = point.X
			b.MaxX = point.X
			b.MinY = point.Y
			b.MaxY = point.Y
		}
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
	return b
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

// gdalAccessor wraps around the golang gdal wrapper
type gdalAccessor struct {
	d *gdal.Dataset
	r *gdal.RasterBand
}

func newGDALAccessor(file string) (gdalAccessor, error) {
	d, err := gdal.Open(file, gdal.ReadOnly)
	if err != nil {
		return gdalAccessor{}, err
	}
	r := d.RasterBand(1)
	return gdalAccessor{
		d: &d,
		r: &r,
	}, nil
}

func (g gdalAccessor) close() {
	g.d.Close()
}

func (g gdalAccessor) calculateElevation(rasterBBox BoundingBox, p *Point) error {
	bufSizeX := int(1)
	bufSizeY := int(1)
	buf := make([]float32, bufSizeX*bufSizeY)
	sizeX := g.r.XSize()
	sizeY := g.r.YSize()
	// https://gdal.org/tutorials/geotransforms_tut.html
	// InvGeoTransform works exactly like GeoTransform, but converts from image coordinate space to goreference space
	igt := g.d.InvGeoTransform()
	row := int(igt[0] + p.X*igt[1] + p.Y*igt[2])
	col := int(igt[3] + p.X*igt[4] + p.Y*igt[5])
	if row < 0 || row > sizeX || col < 0 || col > sizeY {
		return errors.New("Point lies outside Item BoundingBox")
	}
	var err error
	// C++ API
	// https://gdal.org/api/gdalrasterband_cpp.html
	err = g.r.IO(gdal.Read, row, col, bufSizeX, bufSizeY, buf, bufSizeX, bufSizeY, 0, 0)
	if err != nil {
		return err
	}
	v := float64(buf[0])
	p.Elevation = &v
	return nil
}
