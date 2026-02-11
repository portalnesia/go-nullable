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
	"errors"
	"fmt"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Type represents a custom struct that may be null or not
// present in JSON at all.
type Type[D any] struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid string
	Data    D
}

func NewType[T any](data T, presentValid ...bool) Type[T] {
	d := Type[T]{
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
func NewTypePtr[T any](data T, presentValid ...bool) *Type[T] {
	d := NewType[T](data, presentValid...)
	return &d
}

func (d Type[D]) IsPresent() bool {
	return d.Present
}

func (d Type[D]) IsValid() bool {
	return d.Valid
}

func (d Type[D]) GetValue() interface{} {
	return d.Data
}

func (d Type[D]) Datatypes() *datatypes.JSONType[D] {
	if !d.Present || !d.Valid {
		return nil
	}
	dt := datatypes.NewJSONType(d.Data)
	return &dt
}

func (d Type[D]) Ptr() *D {
	if d.Valid {
		return &d.Data
	}
	return nil
}

var (
	_ driver.Valuer       = (*Type[any])(nil)
	_ sql.Scanner         = (*Type[any])(nil)
	_ json.Marshaler      = (*Type[any])(nil)
	_ json.Unmarshaler    = (*Type[any])(nil)
	_ bson.Marshaler      = (*Type[any])(nil)
	_ bson.Unmarshaler    = (*Type[any])(nil)
	_ msgpack.Marshaler   = (*Type[any])(nil)
	_ msgpack.Unmarshaler = (*Type[any])(nil)
)

// Scan implements sql.Scanner interface
func (d *Type[D]) Scan(value interface{}) error {
	d.Present = true
	if value == nil {
		d.Valid = false
		return nil
	}

	var byteData []byte
	switch v := value.(type) {
	case []byte:
		byteData = v
	case string:
		byteData = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}
	d.Valid = true
	return json.Unmarshal(byteData, &d.Data)
}

// Value implements driver.Valuer interface
func (d Type[D]) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	val, err := json.Marshal(d.Data)
	if err != nil {
		return nil, err
	}
	return string(val), nil
}

// MarshalJSON implements json.Marshaler interface.
// Bug: Marshal undefined value
func (d Type[D]) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *Type[D]) UnmarshalJSON(data []byte) error {
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

// MarshalBSON implements bson.Marshaler interface.
func (d Type[D]) MarshalBSON() (byt []byte, err error) {
	var tmp *bool
	_, byt, err = bson.MarshalValue(tmp)
	if !d.Present {
		return byt, err
	} else if !d.Valid {
		return byt, err
	}
	return bson.Marshal(d.Data)
}

// UnmarshalBSON implements bson.Marshaler interface.
func (d *Type[D]) UnmarshalBSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := bson.Unmarshal(data, &d.Data); err != nil {
		return err
	}
	d.Valid = true
	return nil
}

// GormDataType gorm common data type
func (Type[D]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Type[D]) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

// MarshalMsgpack implements msgpack.Marshaler interface.
func (d Type[D]) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *Type[D]) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *D
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

func (Type[D]) FiberConverter(value string) reflect.Value {
	var tmp D
	s := Type[D]{
		true,
		false,
		tmp,
	}

	if err := json.Unmarshal([]byte(value), &tmp); err != nil {
		s = NewType(tmp, true, true)
	}

	return reflect.ValueOf(s)
}
