package yc_lockbox_unpack

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func unpackText(field reflect.Value, value string) error {
	if field.Kind() == reflect.Ptr {
		switch {
		case value == "":
			return nil
		case field.IsNil():
			field.Set(reflect.New(field.Type().Elem()))
		}

		field = field.Elem()
	}

	if !field.CanInterface() {
		return nil
	}

	fieldIf := field.Interface()
	switch fieldIf.(type) {
	case bool:
		value = strings.TrimSpace(value)
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case int, int8, int16, int32, int64:
		value = strings.TrimSpace(value)
		i, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case uint, uint8, uint16, uint32, uint64:
		value = strings.TrimSpace(value)
		u, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case float32, float64:
		value = strings.TrimSpace(value)
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case string:
		field.SetString(value)
	case []byte:
		field.SetBytes([]byte(value))
	default:
		for field.CanAddr() {
			field = field.Addr()
		}

		umarshaller, ok := field.Interface().(encoding.TextUnmarshaler)
		if ok {
			return umarshaller.UnmarshalText([]byte(value))
		}

		return fmt.Errorf("don't know how to parse type: %s", field.Type())
	}

	return nil
}

func unpackBinary(field reflect.Value, value []byte) error {
	if field.Kind() == reflect.Ptr {
		switch {
		case value == nil:
			return nil
		case field.IsNil():
			field.Set(reflect.New(field.Type().Elem()))
		}

		field = field.Elem()
	}

	if !field.CanInterface() {
		return nil
	}

	fieldIf := field.Interface()
	switch fieldIf.(type) {
	case string:
		field.SetString(string(value))
	case []byte:
		field.SetBytes(value)
	default:
		for field.CanAddr() {
			field = field.Addr()
		}

		umarshaller, ok := field.Interface().(encoding.BinaryUnmarshaler)
		if ok {
			return umarshaller.UnmarshalBinary(value)
		}

		return fmt.Errorf("don't know how to parse type: %s", field.Type())
	}

	return nil
}
