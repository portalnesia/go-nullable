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
	"reflect"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

// Float represents a float that may be null or not
// present in JSON at all.
type Float struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid float
	Data    float64
}

func NewFloat(data float64, presentValid ...bool) Float {
	d := Float{
		Data:    data,
		Present: true,
		Valid:   true,
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
func NewFloatPtr(data float64, presentValid ...bool) *Float {
	d := NewFloat(data, presentValid...)
	return &d
}

func (d Float) IsPresent() bool {
	return d.Present
}

func (d Float) IsValid() bool {
	return d.Valid
}

func (d Float) GetValue() interface{} {
	return d.Data
}

func (d Float) Null() null.Float {
	return null.NewFloat(d.Data, d.Present && d.Valid)
}

func (d Float) Ptr() *float64 {
	if d.Valid {
		return &d.Data
	}
	return nil
}

var (
	_ driver.Valuer       = (*Float)(nil)
	_ sql.Scanner         = (*Float)(nil)
	_ json.Marshaler      = (*Float)(nil)
	_ json.Unmarshaler    = (*Float)(nil)
	_ bson.Marshaler      = (*Float)(nil)
	_ bson.Unmarshaler    = (*Float)(nil)
	_ msgpack.Marshaler   = (*Float)(nil)
	_ msgpack.Unmarshaler = (*Float)(nil)
)

// Scan implements sql.Scanner interface
func (d *Float) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullFloat64
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Float64
	return nil
}

// Value implements driver.Valuer interface
func (d Float) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (d Float) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *Float) UnmarshalJSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := json.Unmarshal(data, &d.Data); err != nil {
		return nil
	}

	d.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (d Float) MarshalBSON() (byt []byte, err error) {
	var tmp *float64
	_, byt, err = bson.MarshalValue(tmp)
	if !d.Present {
		return byt, err
	} else if !d.Valid {
		return byt, err
	}
	_, byt, err = bson.MarshalValue(d.Data)
	return byt, err
}

// UnmarshalBSON implements bson.Marshaler interface.
func (d *Float) UnmarshalBSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := bson.Unmarshal(data, &d.Data); err != nil {
		return nil
	}

	d.Valid = true
	return nil
}

// MarshalMsgpack implements msgpack.Marshaler interface.
func (d Float) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *Float) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *float64
	if err := msgpack.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		d.Valid = false
		return nil
	}

	d.Valid = true
	d.Data = *val
	return nil
}

func (Float) FiberConverter(value string) reflect.Value {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		a := NewFloat(f, true, false)
		return reflect.ValueOf(a)
	}
	a := NewFloat(f, true, true)
	return reflect.ValueOf(a)
}
