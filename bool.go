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
	"encoding/json"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.portalnesia.com/utils"

	"gopkg.in/guregu/null.v4"
)

// Bool represents a bool that may be null or not
// present in JSON at all.
type Bool struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid bool
	Data    bool
}

func (d Bool) IsPresent() bool {
	return d.Present
}

func (d Bool) IsValid() bool {
	return d.Valid
}

func (d Bool) GetValue() interface{} {
	return d.Data
}

func NewBool(data bool, presentValid ...bool) Bool {
	d := Bool{
		Present: true,
		Valid:   true,
		Data:    data,
	}

	if len(presentValid) > 0 {
		d.Present = presentValid[0]
		if len(presentValid) > 1 {
			d.Valid = presentValid[1]
		}
	}
	return d
}
func NewBoolPtr(data bool, presentValid ...bool) *Bool {
	d := NewBool(data, presentValid...)
	return &d
}

func (d Bool) Null() null.Bool {
	return null.NewBool(d.Data, d.Present && d.Valid)
}

func (d Bool) Ptr() *bool {
	if d.Valid {
		return &d.Data
	}
	return nil
}

var (
	_ driver.Valuer       = (*Bool)(nil)
	_ sql.Scanner         = (*Bool)(nil)
	_ json.Marshaler      = (*Bool)(nil)
	_ json.Unmarshaler    = (*Bool)(nil)
	_ bson.Marshaler      = (*Bool)(nil)
	_ bson.Unmarshaler    = (*Bool)(nil)
	_ msgpack.Marshaler   = (*Bool)(nil)
	_ msgpack.Unmarshaler = (*Bool)(nil)
)

// Scan implements sql.Scanner interface
func (d *Bool) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullBool
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Bool
	return nil
}

// Value implements driver.Valuer interface
func (d Bool) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (d Bool) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *Bool) UnmarshalJSON(data []byte) error {
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
func (d Bool) MarshalBSON() (byt []byte, err error) {
	var tmp *bool
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
func (d *Bool) UnmarshalBSON(data []byte) error {
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
func (d Bool) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *Bool) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *bool
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

func (Bool) FiberConverter(value string) reflect.Value {
	b := utils.IsTrue(value)
	a := NewBool(b, true, true)
	return reflect.ValueOf(a)
}
