/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package nullable

import (
	"bytes"
	"encoding/json"
	pg "github.com/lib/pq"
	"reflect"
	"testing"
)

type typeStringArrayTest struct {
	Value StringArray `json:"value"`
}

func TestStringArray_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   typeStringArrayTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: typeStringArrayTest{
				Value: StringArray{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: typeStringArrayTest{
				Value: StringArray{
					Data: pg.StringArray{
						"test",
						"string",
						"array",
					},
					Present: true,
					Valid:   true,
				},
			},
			expect: bytes.NewBufferString(`{"value":["test","string","array"]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var byt []byte
			var err error

			if byt, err = json.Marshal(tt.data); err != nil {
				t.Fatalf("unexpected marshaling error: %s", err)
			}

			if !bytes.Equal(byt, tt.expect.Bytes()) {
				t.Errorf("expected value to be %s got %s", tt.expect, byt)
			}
		})
	}
}

func TestStringArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect StringArray
	}{
		{
			name:   "undefined",
			buf:    bytes.NewBufferString(`{}`),
			expect: StringArray{},
		},
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value": null}`),
			expect: StringArray{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":["test","string","array"]}`),
			expect: StringArray{
				Present: true,
				Valid:   true,
				Data: pg.StringArray{
					"test",
					"string",
					"array",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value StringArray `json:"value"`
			}{}

			if err := json.Unmarshal(tt.buf.Bytes(), &str); err != nil {
				t.Fatalf("unexpected unmarshaling error: %s", err)
			}

			got := str.Value
			if got.Present != tt.expect.Present || got.Valid != tt.expect.Valid || !reflect.DeepEqual(got.Data, tt.expect.Data) {
				t.Errorf("expected value to be %#v got %#v", tt.expect, got)
			}
		})
	}
}
