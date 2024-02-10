package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	pg "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type StringArray struct {
	Present bool // Present is true if key is present in json
	Valid   bool // Valid is true if value is not null and valid string
	Data    pg.StringArray
}

func NewStringArray(data pg.StringArray, present bool, valid bool) StringArray {
	return StringArray{Present: present, Valid: valid, Data: data}
}
func NewStringArrayPtr(data pg.StringArray, present bool, valid bool) *StringArray {
	return &StringArray{Present: present, Valid: valid, Data: data}
}

func (d StringArray) Ptr() *pg.StringArray {
	if d.Valid {
		return &d.Data
	}
	return nil
}

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

	if temp != nil {
		d.Valid = true
	}
	d.Data = temp
	return nil
}

// Value sql.Value interface
func (d StringArray) Value() (driver.Value, error) {
	if !d.Valid || len(d.Data) == 0 {
		return nil, nil
	}

	return d.Data.Value()
}

// MarshalJSON implements json.Marshaler interface.
func (d StringArray) MarshalJSON() ([]byte, error) {
	if !d.Present || !d.Valid {
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

// GormDataType gorm common data type
func (StringArray) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (StringArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
