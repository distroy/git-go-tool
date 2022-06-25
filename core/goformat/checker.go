/*
 * Copyright (C) distroy
 */

package goformat

type Checker interface {
	Check(x *Context) error
}

type Config struct {
	FileLine int  `flag:"default:1000; usage:file line limit. disable if <= 0"`
	Import   bool `flag:"default:true; usage:enable/disable check import"`
	Formated bool `flag:"default:true; usage:enable/disable check file formated"`
	Package  bool `flag:"default:true; usage:enable/disable check package name"`

	FuncInputNum               int  `flag:"default:3; usage:func input num limit. disable if <= 0"`
	FuncOutputNum              int  `flag:"default:3; usage:func output num limit. disable if <= 0"`
	FuncNamedOutput            bool `flag:"default:true; usage:check func output param if need be named"`
	FuncInputNumWithoutContext bool `flag:"default:true; usage:func input num limit if without context"`
	FuncOutputNumWithoutError  bool `flag:"default:true; usage:func output num limit if without error"`
	FuncContextFirst           bool `flag:"default:true; usage:context should be the firsr input parameter"`
	FuncErrorLast              bool `flag:"default:true; usage:error should be the last output parameter"`
	FuncContextErrorMatch      bool `flag:"default:false; usage:context and error should both be standard, or both not be"`
}

func BuildChecker(cfg *Config) Checker {
	checkers := make([]Checker, 0, 8)

	checkers = append(checkers, FileLineChecker(cfg.FileLine))
	checkers = append(checkers, PackageChecker(cfg.Package))
	checkers = append(checkers, ImportChecker(cfg.Import))
	checkers = append(checkers, FormatChecker(cfg.Formated))
	checkers = append(checkers, FuncParamsChecker(&FuncParamsConfig{
		InputNum:               cfg.FuncInputNum,
		OutputNum:              cfg.FuncOutputNum,
		NamedOutput:            cfg.FuncNamedOutput,
		InputNumWithoutContext: cfg.FuncInputNumWithoutContext,
		OutputNumWithoutError:  cfg.FuncOutputNumWithoutError,
		ContextFirst:           cfg.FuncContextFirst,
		ErrorLast:              cfg.FuncErrorLast,
		ContextErrorMatch:      cfg.FuncContextErrorMatch,
	}))

	return Checkers(checkers...)
}
