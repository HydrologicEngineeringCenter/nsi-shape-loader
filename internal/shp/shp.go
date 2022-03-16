package shp

import "github.com/jonas-p/go-shp"

func NewShp(src string) (*shp.Reader, error) {
	shpf, err := shp.Open(src)
	defer shpf.Close()
	return shpf, err
}
