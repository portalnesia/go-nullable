/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
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

func NewInt(data int64, present bool, valid bool) Int {
	return Int{
		Present: present,
		Valid:   valid,
		Data:    data,
	}
}
func NewIntPtr(data int64, present bool, valid bool) *Int {
	return &Int{
		Present: present,
		Valid:   valid,
		Data:    data,
	}
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

func (Int) FiberConverter(value string) reflect.Value {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		a := NewInt(i, true, false)
		return reflect.ValueOf(a)
	}
	a := NewInt(i, true, true)
	return reflect.ValueOf(a)
}
