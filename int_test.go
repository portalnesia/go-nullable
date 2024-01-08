/*
Copyright © Portalnesia <support@portalnesia.com>
*/
package nullable

import (
	"bytes"
	"testing"

	"encoding/json"
)

type intJsonTest struct {
	Value Int `json:"value"`
}

func TestInt_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   intJsonTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: intJsonTest{
				Value: Int{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: intJsonTest{
				Value: Int{
					Present: true,
					Valid:   true,
					Data:    5,
				},
			},
			expect: bytes.NewBufferString(`{"value":5}`),
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

func TestInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect Int
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Int{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":1}`),
			expect: Int{
				Present: true,
				Valid:   true,
				Data:    1,
			},
		},
		{
			name:   "empty",
			buf:    bytes.NewBufferString(`null`),
			expect: Int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Int `json:"value"`
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
