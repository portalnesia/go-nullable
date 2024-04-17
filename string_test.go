/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"testing"

	"encoding/json"
)

type stringJsonTest struct {
	Value String `json:"value,omitempty"`
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   stringJsonTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: stringJsonTest{
				Value: String{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: stringJsonTest{
				Value: String{
					Present: true,
					Valid:   true,
					Data:    "test",
				},
			},
			expect: bytes.NewBufferString(`{"value":"test"}`),
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

func TestString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect String
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: String{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":"string"}`),
			expect: String{
				Present: true,
				Valid:   true,
				Data:    "string",
			},
		},
		{
			name:   "empty",
			buf:    bytes.NewBufferString(`null`),
			expect: String{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value String `json:"value"`
			}{}

			if err := json.Unmarshal(tt.buf.Bytes(), &str); err != nil {
				t.Fatalf("unexpected unmarshaling error: %s", err)
			}

			got := str.Value
			if got.Present != tt.expect.Present || got.Valid != tt.expect.Valid || got.Data != tt.expect.Data {
				t.Errorf("expected value to be %#v got %#v", tt.expect, got)
			}
		})
	}
}
