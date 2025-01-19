/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.portalnesia.com/utils"
	"reflect"

	"gopkg.in/guregu/null.v4"
)

// Bool represents a bool that may be null or not
// present in json at all.
type Bool struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid bool
	Data    bool
}

func NewBool(data bool, presentValid ...bool) Bool {
	d := Bool{Present: true, Valid: true, Data: data}

	if len(presentValid) > 0 {
		d.Present = presentValid[0]
		d.Valid = false
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

// sql.Value interface
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

// sql.Value interface
func (d Bool) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (i Bool) MarshalJSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := json.Unmarshal(data, &b.Data); err != nil {
		return nil
	}

	b.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (i Bool) MarshalBSON() (byt []byte, err error) {
	var tmp *bool
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
func (b *Bool) UnmarshalBSON(data []byte) error {
	b.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if err := bson.Unmarshal(data, &b.Data); err != nil {
		return nil
	}

	b.Valid = true
	return nil
}

func (Bool) FiberConverter(value string) reflect.Value {
	b := utils.IsTrue(value)
	a := NewBool(b, true, true)
	return reflect.ValueOf(a)
}
