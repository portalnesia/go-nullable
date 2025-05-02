/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package nullable

import (
	"bytes"
	"testing"

	"encoding/json"
)

type nestedValue struct {
	Nested string `json:"nested"`
}
type testValue struct {
	Data nestedValue `json:"data"`
	//Undefined Type[nestedValue] `json:"undefined,omitempty"`
}

type typeJsonTest struct {
	Value Type[testValue] `json:"value,omitempty"`
}

func TestType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   typeJsonTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: typeJsonTest{
				Value: Type[testValue]{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: typeJsonTest{
				Value: Type[testValue]{
					Present: true,
					Valid:   true,
					Data:    testValue{Data: nestedValue{Nested: "nested value"}},
				},
			},
			expect: bytes.NewBufferString(`{"value":{"data":{"nested":"nested value"}}}`),
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

func TestType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect Type[testValue]
	}{
		{
			name:   "undefined",
			buf:    bytes.NewBufferString(`{}`),
			expect: Type[testValue]{},
		},
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value": null}`),
			expect: Type[testValue]{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":{"data":{"nested":"nested value"}}}`),
			expect: Type[testValue]{
				Present: true,
				Valid:   true,
				Data: testValue{
					Data: nestedValue{
						Nested: "nested value",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Type[testValue] `json:"value"`
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
