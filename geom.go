/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package nullable

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/geojson"
	"go.mongodb.org/mongo-driver/bson"
)

type GeomPoint struct {
	geojson.Point
}

var (
	_ sql.Scanner      = (*GeomPoint)(nil)
	_ json.Unmarshaler = (*GeomPoint)(nil)
	_ bson.Unmarshaler = (*GeomPoint)(nil)
)

// Scan implements sql.Scanner interface
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

// UnmarshalJSON implements json.Marshaler interface.
func (g *GeomPoint) UnmarshalJSON(data []byte) error {
	p := orb.Point{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Point = geojson.Point(p)
	return nil
}

// UnmarshalBSON implements bson.Marshaler interface.
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

var (
	_ sql.Scanner      = (*GeomPolygon)(nil)
	_ json.Unmarshaler = (*GeomPolygon)(nil)
	_ bson.Unmarshaler = (*GeomPolygon)(nil)
)

// Scan implements sql.Scanner interface
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

// UnmarshalJSON implements json.Marshaler interface.
func (g *GeomPolygon) UnmarshalJSON(data []byte) error {
	p := orb.Polygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.Polygon = geojson.Polygon(p)
	return nil
}

// UnmarshalBSON implements bson.Marshaler interface.
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

var (
	_ sql.Scanner      = (*GeomMultiPolygon)(nil)
	_ json.Unmarshaler = (*GeomMultiPolygon)(nil)
	_ bson.Unmarshaler = (*GeomMultiPolygon)(nil)
)

// Scan implements sql.Scanner interface
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

// UnmarshalJSON implements json.Marshaler interface.
func (g *GeomMultiPolygon) UnmarshalJSON(data []byte) error {
	p := orb.MultiPolygon{}
	gs := ewkb.Scanner(&p)
	err := gs.Scan(data)
	if err != nil {
		return err
	}
	g.MultiPolygon = geojson.MultiPolygon(p)
	return nil
}

// UnmarshalBSON implements bson.Marshaler interface.
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
