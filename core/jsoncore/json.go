/*
 * Copyright (C) distroy
 */

package jsoncore

import (
	"bytes"
	"encoding/json"

	"github.com/distroy/git-go-tool/core/strcore"
)

type Any interface{}

func Marshal(v Any) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func MarshalToString(v Any) (string, error) {
	d, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return strcore.BytesToStrUnsafe(d), nil
}

func MustMarshal(v Any) []byte {
	d, err := Marshal(v)
	if err != nil {
		panic(err)
	}
	return d
}

func MustMarshalToString(v Any) string {
	return strcore.BytesToStrUnsafe(MustMarshal(v))
}

func Unmarshal(d []byte, v Any) error {
	return json.Unmarshal(d, v)
}

func UnmarshalFromString(d string, v Any) error {
	return json.Unmarshal(strcore.StrToBytesUnsafe(d), v)
}
