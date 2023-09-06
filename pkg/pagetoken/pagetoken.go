package pagetoken

import (
	"encoding/base64"

	"github.com/fxamacker/cbor"
)

type Marshaler interface {
	MarshalPageToken() (string, error)
}

type Unmarshaler interface {
	UnmarshalPageToken(data string) error
}

func Marshal(v any) (string, error) {
	if w, ok := v.(Marshaler); ok {
		return w.MarshalPageToken()
	}

	data, err := cbor.Marshal(v, cbor.EncOptions{})
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func Unmarshal(pageToken string, v any) error {
	if w, ok := v.(Unmarshaler); ok {
		return w.UnmarshalPageToken(pageToken)
	}

	data, err := base64.URLEncoding.DecodeString(pageToken)
	if err != nil {
		return err
	}

	if err := cbor.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
