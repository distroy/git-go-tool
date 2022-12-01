/*
 * Copyright (C) distroy
 */

package jsoncore

import (
	"bytes"
	"encoding/json"

	"github.com/distroy/git-go-tool/core/strcore"
)

type any interface{}

func Marshal(v any) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	if err := e.Encode(v); err != nil {
		return nil, err
	}

	s := b.Bytes()
	if l := len(s) - 1; l >= 0 && s[l] == '\n' {
		s = s[:l]
	}
	return s, nil
}

func MarshalToString(v any) (string, error) {
	d, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return strcore.BytesToStrUnsafe(d), nil
}

func MustMarshal(v any) []byte {
	d, err := Marshal(v)
	if err != nil {
		panic(err)
	}
	return d
}

func MustMarshalToString(v any) string {
	return strcore.BytesToStrUnsafe(MustMarshal(v))
}

func Unmarshal(d []byte, v any) error {
	return json.Unmarshal(d, v)
}

func UnmarshalFromString(d string, v any) error {
	return json.Unmarshal(strcore.StrToBytesUnsafe(d), v)
}
