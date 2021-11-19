package yc_lockbox_unpack

import (
	"encoding"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/lockbox/v1"
)

func ptrStr(s string) *string {
	return &s
}

func TestUnpackPayloadText(t *testing.T) {
	type testStruct struct {
		AccessToken *string `secretKey:"ACCESS_TOKEN"`
		TokenSecret string  `secretKey:"TOKEN_SECRET"`
		I           int     `secretKey:"i"`
	}

	testPayload := &lockbox.Payload{
		Entries: []*lockbox.Payload_Entry{
			{
				Key:   "ACCESS_TOKEN",
				Value: &lockbox.Payload_Entry_TextValue{TextValue: "much token"},
			},
			{
				Key:   "TOKEN_SECRET",
				Value: &lockbox.Payload_Entry_TextValue{TextValue: "very secret"},
			},
			{
				Key:   "unused key",
				Value: &lockbox.Payload_Entry_TextValue{TextValue: "notice me"},
			},
			{
				Key:   "i",
				Value: &lockbox.Payload_Entry_TextValue{TextValue: "  9832476  "},
			},
		},
	}

	secret := testStruct{}
	err := UnpackPayload(testPayload, &secret)
	if assert.NoError(t, err) {
		expected := testStruct{
			AccessToken: ptrStr("much token"),
			TokenSecret: "very secret",
			I:           9832476,
		}
		assert.Equal(t, expected, secret)
	}
}

func TestUnpackPayloadBinary(t *testing.T) {
	var _ encoding.BinaryUnmarshaler = &url.URL{}

	type testStruct struct {
		Bytes []byte   `secretKey:"byteZ"`
		URL   *url.URL `secretKey:"Url"`
	}

	testPayload := &lockbox.Payload{
		Entries: []*lockbox.Payload_Entry{
			{
				Key:   "byteZ",
				Value: &lockbox.Payload_Entry_BinaryValue{BinaryValue: []byte("yay bytes")},
			},
			{
				Key:   "Url",
				Value: &lockbox.Payload_Entry_BinaryValue{BinaryValue: []byte("http://localhost:8080")},
			},
		},
	}

	secret := testStruct{}
	err := UnpackPayload(testPayload, &secret)
	if assert.NoError(t, err) {
		expected := testStruct{
			Bytes: []byte("yay bytes"),
			URL:   &url.URL{Scheme: "http", Host: "localhost:8080"},
		}
		assert.Equal(t, expected, secret)
	}
}

func TestUnpackPayloadFailType(t *testing.T) {
	err := UnpackPayload(nil, struct{}{})
	if assert.Error(t, err) {
		assert.EqualError(t, err, "v should be pointer to struct, got: struct {}")
	}

	err = UnpackPayload(nil, nil)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "v should be pointer to struct, got: <nil>")
	}

	err = UnpackPayload(nil, new(int))
	if assert.Error(t, err) {
		assert.EqualError(t, err, "v should be pointer to struct, got: *int")
	}
}
