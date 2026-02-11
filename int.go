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

// Int represents an int64 that may be null or not
// present in JSON at all.
type Int struct {
	Present bool // Present is true if key is present in JSON
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

func (d Int) IsPresent() bool {
	return d.Present
}

func (d Int) IsValid() bool {
	return d.Valid
}

func (d Int) GetValue() interface{} {
	return d.Data
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

var (
	_ driver.Valuer       = (*Int)(nil)
	_ sql.Scanner         = (*Int)(nil)
	_ json.Marshaler      = (*Int)(nil)
	_ json.Unmarshaler    = (*Int)(nil)
	_ bson.Marshaler      = (*Int)(nil)
	_ bson.Unmarshaler    = (*Int)(nil)
	_ msgpack.Marshaler   = (*Int)(nil)
	_ msgpack.Unmarshaler = (*Int)(nil)
)

// Scan implements sql.Scanner interface
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

// Value implements driver.Valuer interface
func (d Int) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (d Int) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *Int) UnmarshalJSON(data []byte) error {
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
func (d Int) MarshalBSON() (byt []byte, err error) {
	var tmp *int
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
func (d *Int) UnmarshalBSON(data []byte) error {
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
func (d Int) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *Int) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *int64
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

func (Int) FiberConverter(value string) reflect.Value {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		a := NewInt(i, true, false)
		return reflect.ValueOf(a)
	}
	a := NewInt(i, true, true)
	return reflect.ValueOf(a)
}
