/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strconv"

	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

type Int struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid int64
	Data    int64
}

func NewInt(data int64, presentValid ...bool) Int {
	d := Int{
		Present: true,
		Valid:   true,
		Data:    data,
	}

	if len(presentValid) > 0 {
		d.Present = presentValid[0]
		d.Valid = false
		if len(presentValid) > 1 {
			d.Valid = presentValid[1]
		}
	}
	return d
}

func NewIntPtr(data int64, presentValid ...bool) *Int {
	d := NewInt(data, presentValid...)
	return &d
}

func (d Int) Null() null.Int {
	return null.NewInt(d.Data, d.Present && d.Valid)
}

func (d Int) Ptr() *int64 {
	if d.Valid {
		return &d.Data
	}
	return nil
}

// sql.Value interface
func (d *Int) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Int64
	return nil
}

// sql.Value interface
func (d Int) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
// Bug: Marshal undefined value
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (i *Int) UnmarshalJSON(data []byte) error {
	i.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := json.Unmarshal(data, &i.Data); err != nil {
		return nil
	}

	i.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (i Int) MarshalBSON() (byt []byte, err error) {
	var tmp *int
	_, byt, err = bson.MarshalValue(tmp)
	if !i.Present {
		return byt, err
	} else if !i.Valid {
		return byt, err
	}
	_, byt, err = bson.MarshalValue(i.Data)
	return byt, err
}

// UnmarshalBSON implements bson.Marshaler interface.
func (i *Int) UnmarshalBSON(data []byte) error {
	i.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := bson.Unmarshal(data, &i.Data); err != nil {
		return nil
	}

	i.Valid = true
	return nil
}

func (Int) FiberConverter(value string) reflect.Value {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		a := NewInt(i, true, false)
		return reflect.ValueOf(a)
	}
	a := NewInt(i, true, true)
	return reflect.ValueOf(a)
}
