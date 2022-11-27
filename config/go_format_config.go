/*
 * Copyright (C) distroy
 */

package config

import (
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/core/ptrcore"
)

var DefaultGoFormat = &GoFormatConfig{
	FileLine: ptrcore.NewInt(1000),
	Import:   ptrcore.NewBool(true),
	Formated: ptrcore.NewBool(true),
	Package:  ptrcore.NewBool(true),

	FuncInputNum:               ptrcore.NewInt(3),
	FuncOutputNum:              ptrcore.NewInt(3),
	FuncNamedOutput:            ptrcore.NewBool(true),
	FuncInputNumWithoutContext: ptrcore.NewBool(true),
	FuncOutputNumWithoutError:  ptrcore.NewBool(true),
	FuncContextFirst:           ptrcore.NewBool(true),
	FuncErrorLast:              ptrcore.NewBool(true),
	FuncContextErrorMatch:      ptrcore.NewBool(false),
}

type GoFormatConfig struct {
	FileLine *int  `yaml:"file-line" flag:"default:1000; usage:file line limit. 0=disable"`
	Import   *bool `yaml:"import"    flag:"default:true; usage:enable/disable check import"`
	Formated *bool `yaml:"formated"  flag:"default:true; usage:enable/disable check file formated"`
	Package  *bool `yaml:"package"   flag:"default:true; usage:enable/disable check package name"`

	FuncInputNum               *int  `yaml:"func-input-num"                 flag:"default:3; usage:func input num limit. 0=disable"`
	FuncOutputNum              *int  `yaml:"func-output-num"                flag:"default:3; usage:func output num limit. 0=disable"`
	FuncNamedOutput            *bool `yaml:"func-named-output"              flag:"default:true; usage:check func output param if need be named"`
	FuncInputNumWithoutContext *bool `yaml:"func-input-num-without-context" flag:"default:true; usage:func input num limit if without context"`
	FuncOutputNumWithoutError  *bool `yaml:"func-output-num-without-error"  flag:"default:true; usage:func output num limit if without error"`
	FuncContextFirst           *bool `yaml:"func-context-first"             flag:"default:true; usage:context should be the firsr input parameter"`
	FuncErrorLast              *bool `yaml:"func-error-last"                flag:"default:true; usage:error should be the last output parameter"`
	FuncContextErrorMatch      *bool `yaml:"func-context-error-match"       flag:"bool; usage:context and error should both be standard, or both not be"`
}

func (c *GoFormatConfig) ToConfig() *goformat.Config {
	return &goformat.Config{
		FileLine: *c.FileLine,
		Import:   *c.Import,
		Formated: *c.Formated,
		Package:  *c.Package,

		FuncInputNum:               *c.FuncInputNum,
		FuncOutputNum:              *c.FuncOutputNum,
		FuncNamedOutput:            *c.FuncNamedOutput,
		FuncInputNumWithoutContext: *c.FuncInputNumWithoutContext,
		FuncOutputNumWithoutError:  *c.FuncOutputNumWithoutError,
		FuncContextFirst:           *c.FuncContextFirst,
		FuncErrorLast:              *c.FuncErrorLast,
		FuncContextErrorMatch:      *c.FuncContextErrorMatch,
	}
}
