/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"time"

	"encoding/json"
	"gopkg.in/guregu/null.v4"
)

type Time struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid string
	Data    time.Time
	carbon  carbon.Carbon
}

func NewTime(data time.Time, presentValid ...bool) Time {
	d := Time{Present: true, Valid: true, Data: data}
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

func (d Time) Null() null.Time {
	return null.NewTime(d.Data, d.Present && d.Valid)
}

func (d Time) Ptr() *time.Time {
	if d.Valid {
		return &d.Data
	}
	return nil
}

func (d Time) Carbon() carbon.Carbon {
	return d.carbon
}

// sql.Value interface
func (d *Time) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Time
	d.carbon = carbon.FromStdTime(i.Time)
	return nil
}

// sql.Value interface
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
	d.Data = carbonTime.ToStdTime()
	d.Valid = true
	d.carbon = carbonTime
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (d Time) MarshalBSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return bson.Marshal(d.Data)
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
	d.Data = carbonTime.ToStdTime()
	d.Valid = true
	d.carbon = carbonTime
	return nil
}

func (Time) FiberConverter(value string) reflect.Value {
	c := carbon.Parse(value)
	if c.IsValid() {
		a := NewTime(c.ToStdTime(), true, true)
		return reflect.ValueOf(a)
	} else {
		a := NewTime(time.Now(), true, false)
		return reflect.ValueOf(a)
	}
}
