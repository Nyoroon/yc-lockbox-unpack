package yc_lockbox_unpack

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/lockbox/v1"
)

const tag = "secretKey"

func UnpackPayload(payload *lockbox.Payload, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v should be pointer to struct, got: %v", reflect.TypeOf(v))
	}

	rv = rv.Elem()

	tagToFieldName := make(map[string]string)
	nameToFieldName := make(map[string]string)

	for _, field := range reflect.VisibleFields(rv.Type()) {
		if fieldTag, ok := field.Tag.Lookup(tag); ok {
			tagToFieldName[fieldTag] = field.Name
		} else {
			fieldName := strings.ToLower(field.Name)
			nameToFieldName[fieldName] = field.Name
		}
	}

	for _, entry := range payload.Entries {
		key := entry.GetKey()
		fieldName, ok := tagToFieldName[key]
		if !ok {
			key = strings.ToLower(key)
			fieldName, ok = nameToFieldName[key]
			if !ok {
				continue
			}
		}

		field := rv.FieldByName(fieldName)

		rawValue := entry.GetValue()
		switch rawValue := rawValue.(type) {
		case *lockbox.Payload_Entry_TextValue:
			if err := unpackText(field, rawValue.TextValue); err != nil {
				return fmt.Errorf("unable to parse text entry %s (field %s): %w", key, fieldName, err)
			}
		case *lockbox.Payload_Entry_BinaryValue:
			if err := unpackBinary(field, rawValue.BinaryValue); err != nil {
				return fmt.Errorf("unable to parse binary entry %s (field %s): %w", key, fieldName, err)
			}
		default:
			return fmt.Errorf("unknown payload value type: %s", reflect.TypeOf(rawValue))
		}
	}

	return nil
}
