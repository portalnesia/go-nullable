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

type floatJsonTest struct {
	Value Float `json:"value"`
}

func TestFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   floatJsonTest
		expect *bytes.Buffer
	}{
		{
			name: "null value",
			data: floatJsonTest{
				Value: Float{
					Present: true,
					Valid:   false,
				},
			},
			expect: bytes.NewBufferString(`{"value":null}`),
		},
		{
			name: "valid value",
			data: floatJsonTest{
				Value: Float{
					Present: true,
					Valid:   true,
					Data:    0.5,
				},
			},
			expect: bytes.NewBufferString(`{"value":0.5}`),
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

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		expect Float
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Float{
				Present: true,
			},
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":1.1}`),
			expect: Float{
				Present: true,
				Valid:   true,
				Data:    1.1,
			},
		},
		{
			name:   "empty",
			buf:    bytes.NewBufferString(`null`),
			expect: Float{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Float `json:"value"`
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
