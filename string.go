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

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

// String represents a string that may be null or not
// present in JSON at all.
type String struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid string
	Data    string
}

func NewString(data string, presentValid ...bool) String {
	d := String{
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

func NewStringPtr(data string, presentValid ...bool) *String {
	d := NewString(data, presentValid...)
	return &d
}

func (d String) IsPresent() bool {
	return d.Present
}

func (d String) IsValid() bool {
	return d.Valid
}

func (d String) GetValue() interface{} {
	return d.Data
}

func (d String) Null() null.String {
	return null.NewString(d.Data, d.Present && d.Valid && d.Data != "")
}

func (d String) Ptr() *string {
	if d.Valid {
		return &d.Data
	}
	return nil
}

var (
	_ driver.Valuer       = (*String)(nil)
	_ sql.Scanner         = (*String)(nil)
	_ json.Marshaler      = (*String)(nil)
	_ json.Unmarshaler    = (*String)(nil)
	_ bson.Marshaler      = (*String)(nil)
	_ bson.Unmarshaler    = (*String)(nil)
	_ msgpack.Marshaler   = (*String)(nil)
	_ msgpack.Unmarshaler = (*String)(nil)
)

// Scan implements sql.Scanner interface
func (d *String) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.String
	return nil
}

// Value implements driver.Valuer interface
func (d String) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (d String) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *String) UnmarshalJSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}

	if err := json.Unmarshal(data, &d.Data); err != nil {
		return nil
	}

	d.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (d String) MarshalBSON() (byt []byte, err error) {
	var tmp *string
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
func (d *String) UnmarshalBSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}

	if err := bson.Unmarshal(data, &d.Data); err != nil {
		return nil
	}

	d.Valid = true
	return nil
}

// MarshalMsgpack implements msgpack.Marshaler interface.
func (d String) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *String) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *string
	if err := msgpack.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		d.Valid = false
		return nil
	}

	d.Valid = len(*val) > 0
	d.Data = *val
	return nil
}

func (String) FiberConverter(value string) reflect.Value {
	a := NewString(value, true, true)
	return reflect.ValueOf(a)
}
