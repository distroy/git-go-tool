/*
 * Copyright (C) distroy
 */

package goformat

type Error = error

type Checker interface {
	Check(x *Context) Error
}

type Config struct {
	FileLine int
	Import   bool
	Formated bool
	Package  bool

	JsonLabel bool

	FuncInputNum               int
	FuncOutputNum              int
	FuncNamedOutput            bool
	FuncInputNumWithoutContext bool
	FuncOutputNumWithoutError  bool
	FuncContextFirst           bool
	FuncErrorLast              bool
	FuncContextErrorMatch      bool
}

func BuildChecker(cfg *Config) Checker {
	checkers := make([]Checker, 0, 8)

	checkers = append(checkers, FileLineChecker(cfg.FileLine))
	checkers = append(checkers, PackageChecker(cfg.Package))
	checkers = append(checkers, ImportChecker(cfg.Import))
	checkers = append(checkers, FormatChecker(cfg.Formated))
	checkers = append(checkers, JsonLabelChecker(cfg.JsonLabel))
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
