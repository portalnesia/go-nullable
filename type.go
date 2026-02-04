/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Type[D any] struct {
	Present bool // Present is true if key is present in json
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

func (d *Type[D]) Scan(value interface{}) error {
	d.Present = true
	if value == nil {
		d.Valid = false
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}
	d.Valid = true
	return json.Unmarshal(bytes, &d.Data)
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
func (i Type[D]) MarshalJSON() ([]byte, error) {
	if !i.Present {
		return []byte(`null`), nil
	} else if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (i *Type[D]) UnmarshalJSON(data []byte) error {
	i.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(data, &i.Data); err != nil {
		return err
	}
	i.Valid = true
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (i Type[D]) MarshalBSON() (byt []byte, err error) {
	var tmp *bool
	_, byt, err = bson.MarshalValue(tmp)
	if !i.Present {
		return byt, err
	} else if !i.Valid {
		return byt, err
	}
	return bson.Marshal(i.Data)
}

// UnmarshalBSON implements bson.Marshaler interface.
func (i *Type[D]) UnmarshalBSON(data []byte) error {
	i.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := bson.Unmarshal(data, &i.Data); err != nil {
		return err
	}
	i.Valid = true
	return nil
}

// GormDataType gorm common data type
func (Type[D]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Type[D]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

//func (js Type[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
//	if !js.Valid {
//		return clause.Expr{SQL: "?", Vars: []any{"NULL"}, WithoutParentheses: true}
//	}
//
//	data, _ := js.MarshalJSON()
//	switch db.Dialector.Name() {
//	case "mysql":
//		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
//			return gorm.Expr("CAST(? AS JSON)", string(data))
//		}
//	}
//
//	return gorm.Expr("?", string(data))
//}
