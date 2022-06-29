/*
 * Copyright (C) distroy
 */

package goformat

func UnitTestChecker() Checker {
	return unitTestChecker{}
}

type unitTestChecker struct {
}

func (c unitTestChecker) Check(x *Context) Error {
	return nil
}
