/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/distroy/git-go-tool/core/tagcore"
)

type fieldTags struct {
	Tags    tagcore.Tags
	Name    string
	Default string
	Usage   string
	IsArgs  bool
}

func parseModel(cmd *flag.FlagSet, val reflect.Value) (args reflect.Value) {
	typ := val.Type()
	if typ.Kind() != reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("input flags must be pointer to struct. %s", typ.String()))
	}

	val = val.Elem()

	return parseStructModel(cmd, val)
}

func parseStructModel(cmd *flag.FlagSet, val reflect.Value) (args reflect.Value) {
	typ := val.Type()

	// log.Printf(" === %s", typ.String())

	for i, l := 0, typ.NumField(); i < l; i++ {
		field := typ.Field(i)
		fVal := val.Field(i)

		if v, ok := fVal.Interface().(flag.Value); ok {
			m := parseFieldTags(field)
			if m == nil {
				// log.Printf(" === aaa %s: %v", fVal.Type().String(), fVal.Interface())
				continue
			}

			// log.Printf(" === %s: %v", fVal.Type().String(), fVal.Interface())
			cmd.Var(v, m.Name, m.Usage)
			continue
		}

		m, tmp := parseFieldModel(cmd, fVal, field)
		if m != nil {
			fillFieldValue(cmd, fVal, m)
		}
		if !args.IsValid() {
			args = tmp
		}
	}

	return args
}

func parseFieldModel(cmd *flag.FlagSet, val reflect.Value, f reflect.StructField) (*fieldTags, reflect.Value) {
	typ := f.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if val.IsNil() {
			val.Set(reflect.New(typ))
		}
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		args := parseStructModel(cmd, val)
		return nil, args
	}

	m := parseFieldTags(f)
	if m == nil {
		return nil, reflect.Value{}
	}

	if m.IsArgs && typ == typeStrings {
		return nil, val
	}

	return m, reflect.Value{}
}

func parseFieldTags(f reflect.StructField) *fieldTags {
	tag, ok := f.Tag.Lookup(tagName)
	if !ok || len(tag) == 0 {
		return &fieldTags{
			Tags: tagcore.New(),
			Name: parseFlagName(f),
		}
	}

	tags := tagcore.Parse(tag)
	if tags.Has("-") {
		return nil
	}
	// log.Printf(" === %s %#v", tag, tags)

	m := &fieldTags{
		Tags:    tags,
		Name:    tags.Get("name"),
		Usage:   tags.Get("usage"),
		Default: tags.Get("default"),
		IsArgs:  tags.Has("args"),
	}

	if len(m.Name) == 0 {
		m.Name = parseFlagName(f)
	}
	// log.Printf(" === 1111 %#v", m)

	return m
}

func fillFieldValue(cmd *flag.FlagSet, val reflect.Value, m *fieldTags) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	// log.Printf(" === 2222 %#v", m)

	fn := fillFlagFuncMap[typ]
	if fn != nil {
		fn(cmd, val, m)
		// log.Printf(" === %s: %v", typ.String(), val.Interface())
	}
	// log.Printf(" === %s", typ.String())
}

func parseFlagName(f reflect.StructField) string {
	name := f.Name
	name = splitStringWord(name, '-')
	name = strings.ToLower(name)
	return name
}

func splitStringWord(s string, sep rune) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	res := make([]rune, 0, len(runes)*2)
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		if !unicode.IsUpper(curr) {
			res = append(res, curr)
			continue
		}

		if i > 0 {
			res = append(res, '-')
		}

		last := curr
		j := i + 1
		for ; j < len(runes); j++ {
			curr := runes[j]
			if unicode.IsUpper(curr) {
				res = append(res, last)
				last = curr
				continue
			}

			if j > i+1 {
				res = append(res, '-')
			}
			res = append(res, last)
			last = curr
			break
		}
		res = append(res, last)
		i = j
	}
	return string(res)
}
