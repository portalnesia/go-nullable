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
	"errors"
	"reflect"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

// Time represents go time that may be null or not
// present in JSON at all.
type Time struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid string
	Data    time.Time
	carbon  *carbon.Carbon
}

func NewTime(data time.Time, presentValid ...bool) Time {
	d := Time{
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
func NewTimePtr(data time.Time, presentValid ...bool) *Time {
	d := NewTime(data, presentValid...)
	return &d
}

func (d Time) IsPresent() bool {
	return d.Present
}

func (d Time) IsValid() bool {
	return d.Valid
}

func (d Time) GetValue() interface{} {
	return d.Data
}

func (d Time) Null() null.Time {
	return null.NewTime(d.Data, d.Present && d.Valid)
}

func (d Time) Ptr() *time.Time {
	if d.Valid {
		return &d.Data
	}
	return nil
}

func (d Time) Carbon() *carbon.Carbon {
	return d.carbon
}

var (
	_ driver.Valuer       = (*Time)(nil)
	_ sql.Scanner         = (*Time)(nil)
	_ json.Marshaler      = (*Time)(nil)
	_ json.Unmarshaler    = (*Time)(nil)
	_ bson.Marshaler      = (*Time)(nil)
	_ bson.Unmarshaler    = (*Time)(nil)
	_ msgpack.Marshaler   = (*Time)(nil)
	_ msgpack.Unmarshaler = (*Time)(nil)
)

// Scan implements sql.Scanner interface
func (d *Time) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Time
	d.carbon = carbon.CreateFromStdTime(i.Time)
	return nil
}

// Value implements driver.Valuer interface
func (d Time) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (d Time) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *Time) UnmarshalJSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	var timeString string

	if err := json.Unmarshal(data, &timeString); err != nil {
		return err
	}

	carbonTime := carbon.Parse(timeString)
	if !carbonTime.IsValid() {
		return errors.New("invalid date string")
	}
	d.Data = carbonTime.StdTime()
	d.Valid = true
	d.carbon = carbonTime
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (d Time) MarshalBSON() (byt []byte, err error) {
	var tmp *time.Time
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
func (d *Time) UnmarshalBSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	var timeString string

	if err := bson.Unmarshal(data, &timeString); err != nil {
		return err
	}

	carbonTime := carbon.Parse(timeString)
	if !carbonTime.IsValid() {
		return errors.New("invalid date string")
	}
	d.Data = carbonTime.StdTime()
	d.Valid = true
	d.carbon = carbonTime
	return nil
}

// MarshalMsgpack implements msgpack.Marshaler interface.
func (d Time) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *Time) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *string
	if err := msgpack.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		d.Valid = false
		return nil
	}

	carbonTime := carbon.Parse(*val)
	if !carbonTime.IsValid() {
		return errors.New("invalid date string")
	}
	d.Data = carbonTime.StdTime()
	d.Valid = true
	d.carbon = carbonTime
	return nil
}

func (Time) FiberConverter(value string) reflect.Value {
	c := carbon.Parse(value)
	a := NewTime(c.StdTime(), true, false)
	if c.IsValid() {
		a.Valid = true
	}

	return reflect.ValueOf(a)
}
