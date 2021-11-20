package yc_lockbox_unpack

import (
	"net"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackText(t *testing.T) {
	tests := []struct {
		typ      reflect.Type
		value    string
		expected interface{}
	}{
		{
			typ:      reflect.TypeOf(true),
			value:    "True",
			expected: true,
		},
		{
			typ:      reflect.TypeOf(1),
			value:    " 1_234_567",
			expected: 1234567,
		},
		{
			typ:      reflect.TypeOf(uint(1)),
			value:    "1_234_567 ",
			expected: uint(1234567),
		},
		{
			typ:      reflect.TypeOf(float32(1.0)),
			value:    "1.23456",
			expected: float32(1.23456),
		},
		{
			typ:      reflect.TypeOf(net.IP{}),
			value:    "192.0.2.1",
			expected: net.IPv4(192, 0, 2, 1),
		},
		{
			typ:      reflect.TypeOf(ptrStr("string")),
			value:    "stringster",
			expected: ptrStr("stringster"),
		},
	}

	var err error
	var value reflect.Value

	for _, test := range tests {
		t.Run(test.typ.String(), func(t *testing.T) {
			value = reflect.New(test.typ).Elem()
			err = unpackText(value, test.value)
			if assert.NoError(t, err) {
				expected := reflect.ValueOf(test.expected)
				assert.Equal(t, expected.Type(), value.Type())
				assert.Equal(t, expected.Interface(), value.Interface())
			}
		})
	}
}

func TestUnpackBinary(t *testing.T) {
	tests := []struct {
		typ      reflect.Type
		value    []byte
		expected interface{}
	}{
		{
			typ:      reflect.TypeOf([]byte{}),
			value:    []byte("bytes"),
			expected: []byte("bytes"),
		},
		{
			typ:      reflect.TypeOf(""),
			value:    []byte("string"),
			expected: "string",
		},
		{
			typ:      reflect.TypeOf(&url.URL{}),
			value:    []byte("http://localhost:8080"),
			expected: &url.URL{Scheme: "http", Host: "localhost:8080"},
		},
	}

	var err error
	var value reflect.Value

	for _, test := range tests {
		t.Run(test.typ.String(), func(t *testing.T) {
			value = reflect.New(test.typ).Elem()
			err = unpackBinary(value, test.value)
			if assert.NoError(t, err) {
				expected := reflect.ValueOf(test.expected)
				assert.Equal(t, expected.Type(), value.Type())
				assert.Equal(t, expected.Interface(), value.Interface())
			}
		})
	}
}
