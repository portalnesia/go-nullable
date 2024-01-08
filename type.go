/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Type[D any] struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid string
	Data    D
}

func NewType[T any](data T, present bool, valid bool) Type[T] {
	return Type[T]{Present: present, Valid: valid, Data: data}
}
func NewTypePtr[T any](data T, present bool, valid bool) *Type[T] {
	return &Type[T]{Present: present, Valid: valid, Data: data}
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

// Value sql.Value interface
func (d Type[D]) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return json.Marshal(d.Data)
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

func (js Type[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := js.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
