/*
 * Copyright (C) distroy
 */

package tagcore

import (
	"strings"
)

type Tags map[string][]string

func New(size ...int) Tags {
	if len(size) > 0 && size[0] > 0 {
		return make(Tags, size[0])
	}
	return make(Tags)
}

func (m Tags) Add(key, value string) {
	key = strings.ToLower(key)
	m[key] = append(m[key], value)
}

func (m Tags) Set(key, value string) {
	key = strings.ToLower(key)
	m[key] = []string{value}
}

func (m Tags) Has(key string) bool {
	key = strings.ToLower(key)

	_, ok := m[key]
	return ok
}

func (m Tags) Values(key string) []string {
	key = strings.ToLower(key)
	return m[key]
}

func (m Tags) Get(key string) string {
	key = strings.ToLower(key)

	v := m[key]
	if len(v) != 0 {
		return v[0]
	}

	return ""
}

func Parse(tag string) Tags {
	tagList := strings.Split(tag, ";")
	m := New()
	for _, v := range tagList {
		if len(v) == 0 {
			continue
		}

		l := strings.SplitN(v, ":", 2)
		k := strings.TrimSpace(l[0])
		if len(k) == 0 {
			continue
		}

		v := ""
		if len(l) >= 2 {
			v = l[1]
		}

		m.Add(k, v)
	}

	return m
}
