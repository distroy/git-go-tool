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
	if !x.IsGoTest() {
		return nil
	}

	return nil
}
