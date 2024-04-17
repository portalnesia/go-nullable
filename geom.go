package nullable

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/geojson"
)

type GeomPoint struct {
	geojson.Point
}

func (g *GeomPoint) Scan(input interface{}) error {
	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("invalid geometry type: get %v", v)
	}
	p := orb.Point{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Point = geojson.Point(p)
	return nil
}

func (g *GeomPoint) UnmarshalBSON(data []byte) error {
	p := orb.Point{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Point = geojson.Point(p)
	return nil
}

type GeomPolygon struct {
	geojson.Polygon
}

func (g *GeomPolygon) Scan(input interface{}) error {
	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("invalid geometry type: get %v", v)
	}
	p := orb.Polygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Polygon = geojson.Polygon(p)
	return nil
}

func (g *GeomPolygon) UnmarshalBSON(data []byte) error {
	p := orb.Polygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Polygon = geojson.Polygon(p)
	return nil
}

type GeomMultiPolygon struct {
	geojson.MultiPolygon
}

func (g *GeomMultiPolygon) Scan(input interface{}) error {
	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("invalid geometry type: get %v", v)
	}
	p := orb.MultiPolygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.MultiPolygon = geojson.MultiPolygon(p)
	return nil
}

func (g *GeomMultiPolygon) UnmarshalBSON(data []byte) error {
	p := orb.MultiPolygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.MultiPolygon = geojson.MultiPolygon(p)
	return nil
}
