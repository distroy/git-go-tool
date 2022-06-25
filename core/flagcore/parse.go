/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"flag"
	"os"
	"reflect"
)

const (
	tagName = "flag"
)

func Parse(v interface{}) {
	ParseArgs(v, os.Args[1:])
}

func ParseArgs(v interface{}, args []string) {
	cmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	val := reflect.ValueOf(v)
	res := parseModel(cmd, val)

	cmd.Parse(args)
	if res.IsValid() {
		res.Set(reflect.ValueOf(cmd.Args()))
	}
}
