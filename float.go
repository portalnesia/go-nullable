/*
Copyright © Portalnesia <support@portalnesia.com>
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

// Float represents a float that may be null or not
// present in json at all.
type Float struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid float
	Data    float64
}

func NewFloat(data float64, presentValid ...bool) Float {
	d := Float{Data: data, Present: true, Valid: true}

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

func (d Float) Null() null.Float {
	return null.NewFloat(d.Data, d.Present && d.Valid)
}

func (d Float) Ptr() *float64 {
	if d.Valid {
		return &d.Data
	}
	return nil
}

// sql.Value interface
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

// sql.Value interface
func (d Float) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (i Float) MarshalJSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (f *Float) UnmarshalJSON(data []byte) error {
	f.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := json.Unmarshal(data, &f.Data); err != nil {
		return nil
	}

	f.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (i Float) MarshalBSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	_, byt, err := bson.MarshalValue(i.Data)
	return byt, err
}

// UnmarshalBSON implements bson.Marshaler interface.
func (f *Float) UnmarshalBSON(data []byte) error {
	f.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := bson.Unmarshal(data, &f.Data); err != nil {
		return nil
	}

	f.Valid = true
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
