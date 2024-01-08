/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"reflect"

	"encoding/json"

	"gopkg.in/guregu/null.v4"
)

// String represents a string that may be null or not
// present in json at all.
type String struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid string
	Data    string
}

func NewString(data string, present bool, valid bool) String {
	return String{Present: present, Valid: valid, Data: data}
}

func NewStringPtr(data string, present bool, valid bool) *String {
	return &String{Present: present, Valid: valid, Data: data}
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

// sql.Value interface
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

// sql.Value interface
func (d String) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (i String) MarshalJSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (s *String) UnmarshalJSON(data []byte) error {
	s.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}

	if err := json.Unmarshal(data, &s.Data); err != nil {
		return nil
	}

	s.Valid = true
	return nil
}

func (String) FiberConverter(value string) reflect.Value {
	a := NewString(value, true, true)
	return reflect.ValueOf(a)
}
