/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"github.com/golang-module/carbon"
	"reflect"
	"time"

	"encoding/json"
	"gopkg.in/guregu/null.v4"
)

type Time struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid string
	Data    time.Time
}

func NewTime(data time.Time, present bool, value bool) Time {
	return Time{Present: present, Valid: value, Data: data}
}
func NewTimePtr(data time.Time, present bool, value bool) *Time {
	return &Time{Present: present, Valid: value, Data: data}
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

// sql.Value interface
func (d *Time) Scan(value interface{}) error {
	d.Present = true

	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}
	d.Valid = i.Valid
	d.Data = i.Time
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

	if err := json.Unmarshal(data, &d.Data); err != nil {
		return err
	}

	d.Valid = true
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
