package elevation

import (
	"errors"
	"fmt"
	"math"

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

func (b BoundingBox) QueryNationalMap() (QueryResult, error) {
	query := NewNationalMapQuery()
	query.setParam("bbox", fmt.Sprintf("%f,%f,%f,%f", b.MinX, b.MinY, b.MaxX, b.MaxY))
	query.setParam("prodFormats", "GeoTIFF")
	r, err := query.sendRequest()
	return r, err
}

// gdalAccessor wraps around the golang gdal wrapper
type gdalAccessor struct {
	d *gdal.Dataset
	r gdal.RasterBand
	a gdal.RasterAttributeTable
}

func newGDALAccessor(file string) (gdalAccessor, error) {
	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return gdalAccessor{}, err
	}
	size := 256
	d := driver.Create(file, size, size, 1, gdal.Byte, nil)
	defer d.Close()
	r := d.RasterBand(1)
	a := r.GetDefaultRAT()
	return gdalAccessor{
		d: &d,
		r: r,
		a: a,
	}, nil
}

func (g gdalAccessor) calculateElevation(rasterBBox BoundingBox, p Point) error {
	// ptr := unsafe.Pointer(new([1 * 1]byte))
	// size := 256
	// var buf float64
	// err := g.r.IO(gdal.Read, 5, 5, size, size, buf, size, size, 0, 0)
	// if err != nil {
	// 	return err
	// }
	buf := []uint8{255}
	sizeX := g.r.XSize()
	sizeY := g.r.YSize()
	pixelSizeX := (rasterBBox.MaxX - rasterBBox.MinX) / float64(sizeX)
	pixelSizeY := (rasterBBox.MaxY - rasterBBox.MinY) / float64(sizeY)
	row := (p.X - rasterBBox.MinX) / pixelSizeX
	col := (p.Y - rasterBBox.MinY) / pixelSizeY
	if int(row) > sizeX || int(col) > sizeY {
		return errors.New("Point lies outside Item BoundingBox")
	}
	// v := g.r.value(int(row), int(col))
	// log.Print(v)
	err := g.r.IO(gdal.Read, int(math.Round(row)), int(math.Round(col)), 1, 1, buf, 1, 1, 0, 0)
	// err := g.r.ReadBlock(int(math.Round(row)), int(math.Round(col)), ptr)
	// p.Elevation = (*float64)(ptr)
	// p.Elevation = &v
	return err
}
