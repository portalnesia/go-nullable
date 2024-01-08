/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"testing"

	"encoding/json"
)

type boolJsonTest struct {
	Value Bool `json:"value"`
}

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   boolJsonTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: boolJsonTest{
				Value: Bool{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: boolJsonTest{
				Value: Bool{
					Present: true,
					Valid:   true,
					Data:    true,
				},
			},
			expect: bytes.NewBufferString(`{"value":true}`),
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
func TestBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect Bool
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Bool{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":true}`),
			expect: Bool{
				Present: true,
				Valid:   true,
				Data:    true,
			},
		},
		{
			name:   "empty",
			buf:    bytes.NewBufferString(`null`),
			expect: Bool{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Bool `json:"value"`
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
