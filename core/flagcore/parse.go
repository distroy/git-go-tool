/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"
)

const (
	tagName = "flag"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func MustParse(v interface{}, args ...[]string) {
	s := NewFlagSet()
	s.Model(v)
	s.MustParse(args...)
}

func Parse(v interface{}, args ...[]string) error {
	s := NewFlagSet()
	s.Model(v)
	return s.Parse(args...)
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

func unquoteUsage(f *Flag) (meta string, usage string) {
	usage = f.Usage
	meta = f.Meta
	if meta != "" {
		if !strings.HasPrefix(meta, "<") && !strings.HasPrefix(meta, ">") {
			meta = fmt.Sprintf("<%s>", meta)
		}

		return meta, usage
	}

	// Look for a back-quoted name, but avoid the strings package.
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					meta = usage[i+1 : j]
					usage = usage[:i] + meta + usage[j+1:]
					return meta, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}

	// No explicit name, so use type if we can find one.
	meta = "<value>"
	switch f.Value.(type) {
	case *boolFlag:
		meta = ""
	case *boolValue:
		meta = "<bool>"
	case *durationValue:
		meta = "<duration>"
	case *float32Value, *float64Value, *float32sValue, *float64sValue:
		meta = "<float>"
	case *intValue, *int64Value, *intsValue, *int64sValue:
		meta = "<int>"
	case *stringValue, *stringsValue:
		meta = "<string>"
	case *uintValue, *uint64Value, *uintsValue, *uint64sValue:
		meta = "<uint>"
	}
	return meta, usage
}

func isZeroValue(value Value, defValue string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return defValue == z.Interface().(Value).String()
}
