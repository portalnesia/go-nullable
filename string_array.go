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
	"reflect"

	pg "github.com/lib/pq"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/schema"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

// StringArray represents an array of string that may be null or not
// present in JSON at all.
//
// When using with bun ORM, do NOT use `type:text[]` or `array` tag on the field,
// as pgdialect will override the appender with arrayAppender which does not support struct types.
// Instead, use nullzero tag and let the driver.Valuer handle the conversion:
//
//	type Model struct {
//	    Tags StringArray `bun:"tags,nullzero"`
//	}
type StringArray struct {
	Present bool // Present is true if key is present in JSON
	Valid   bool // Valid is true if value is not null and valid string
	Data    pg.StringArray
}

func NewStringArray(data pg.StringArray, presentValid ...bool) StringArray {
	d := StringArray{
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
func NewStringArrayPtr(data pg.StringArray, presentValid ...bool) *StringArray {
	d := NewStringArray(data, presentValid...)
	return &d
}

func (d StringArray) Ptr() *pg.StringArray {
	if d.Valid {
		return &d.Data
	}
	return nil
}

func (d StringArray) IsPresent() bool {
	return d.Present
}

func (d StringArray) IsValid() bool {
	return d.Valid
}

func (d StringArray) GetValue() interface{} {
	return d.Data
}

var (
	_ driver.Valuer        = (*StringArray)(nil)
	_ sql.Scanner          = (*StringArray)(nil)
	_ json.Marshaler       = (*StringArray)(nil)
	_ json.Unmarshaler     = (*StringArray)(nil)
	_ bson.Marshaler       = (*StringArray)(nil)
	_ bson.Unmarshaler     = (*StringArray)(nil)
	_ msgpack.Marshaler    = (*StringArray)(nil)
	_ msgpack.Unmarshaler  = (*StringArray)(nil)
	_ schema.QueryAppender = (*StringArray)(nil)
)

// Scan implements sql.Scanner interface
func (d *StringArray) Scan(value interface{}) error {
	d.Present = true
	if value == nil {
		d.Valid = false
		return nil
	}

	var temp pg.StringArray
	if err := temp.Scan(value); err != nil {
		d.Valid = false
		return err
	}

	d.Valid = true
	d.Data = temp
	return nil
}

// Value implements driver.Valuer interface
func (d StringArray) Value() (driver.Value, error) {
	if !d.Valid || len(d.Data) == 0 {
		return nil, nil
	}

	return d.Data.Value()
}

// AppendQuery implements bun/dialect.AppendQuery interface for PostgreSQL array support.
func (d StringArray) AppendQuery(gen schema.QueryGen, b []byte) ([]byte, error) {
	if !d.Valid {
		return dialect.AppendNull(b), nil
	}
	return pgdialect.Array([]string(d.Data)).AppendQuery(gen, b)
}

// MarshalJSON implements json.Marshaler interface.
func (d StringArray) MarshalJSON() ([]byte, error) {
	if !d.Present {
		return []byte(`null`), nil
	} else if !d.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(d.Data)
}

// UnmarshalJSON implements json.Marshaler interface.
func (d *StringArray) UnmarshalJSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(data, &d.Data); err != nil {
		return err
	}
	if len(d.Data) > 0 {
		d.Valid = true
	}
	return nil
}

// MarshalBSON implements bson.Marshaler interface.
func (d StringArray) MarshalBSON() (byt []byte, err error) {
	var tmp *string
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
func (d *StringArray) UnmarshalBSON(data []byte) error {
	d.Present = true

	if bytes.Equal(data, []byte("null")) {
		return nil
	}
	if err := bson.Unmarshal(data, &d.Data); err != nil {
		return err
	}
	if len(d.Data) > 0 {
		d.Valid = true
	}
	return nil
}

// GormDataType gorm common data type
func (StringArray) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (StringArray) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "text[]"
	}
	return ""
}

// MarshalMsgpack implements msgpack.Marshaler interface.
func (d StringArray) MarshalMsgpack() ([]byte, error) {
	if !d.Present || !d.Valid {
		return msgpack.Marshal(nil)
	}
	return msgpack.Marshal(d.Data)
}

// UnmarshalMsgpack implements msgpack.Unmarshaler interface.
func (d *StringArray) UnmarshalMsgpack(data []byte) error {
	d.Present = true // Jika fungsi ini dipanggil, berarti key-nya ada di payload

	var val *[]string
	if err := msgpack.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		d.Valid = false
		return nil
	}

	d.Valid = len(*val) > 0
	d.Data = *val
	return nil
}

func (StringArray) FiberConverter(value string) reflect.Value {
	var tmp pg.StringArray
	s := StringArray{
		true,
		false,
		pg.StringArray{},
	}

	if err := json.Unmarshal([]byte(value), &tmp); err != nil {
		s = NewStringArray(tmp, true, true)
	}

	return reflect.ValueOf(s)
}

//func (d StringArray) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
//	data, _ := d.MarshalJSON()
//	switch db.Dialector.Name() {
//	case "mysql":
//		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
//			return gorm.Expr("CAST(? AS JSON)", string(data))
//		}
//	}
//
//	return clause.Expr{}
//}
