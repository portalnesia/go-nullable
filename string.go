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
	"go.mongodb.org/mongo-driver/bson"
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

func (s String) Null() null.String {
	return null.NewString(s.Data, s.Present && s.Valid && s.Data != "")
}

func (s String) Ptr() *string {
	if s.Valid {
		return &s.Data
	}
	return nil
}

// sql.Value interface
func (s *String) Scan(value interface{}) error {
	s.Present = true

	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	s.Valid = i.Valid
	s.Data = i.String
	return nil
}

// sql.Value interface
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Data, nil
}

// MarshalJSON implements json.Marshaler interface.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Present {
		return []byte(`null`), nil
	} else if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Data)
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

// MarshalBSON implements bson.Marshaler interface.
func (s String) MarshalBSON() (byt []byte, err error) {
	var tmp *string
	_, byt, err = bson.MarshalValue(tmp)
	if !s.Present {
		return byt, err
	} else if !s.Valid {
		return byt, err
	}
	_, byt, err = bson.MarshalValue(s.Data)
	return byt, err
}

// UnmarshalBSON implements bson.Marshaler interface.
func (s *String) UnmarshalBSON(data []byte) error {
	s.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}

	if err := bson.Unmarshal(data, &s.Data); err != nil {
		return nil
	}

	s.Valid = true
	return nil
}

func (String) FiberConverter(value string) reflect.Value {
	a := NewString(value, true, true)
	return reflect.ValueOf(a)
}
