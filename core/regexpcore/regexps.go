/*
 * Copyright (C) distroy
 */

package regexpcore

import (
	"encoding/json"
	"regexp"
)

var DefaultExcludes = []string{
	`^vendor/`,
	`/vendor/`,
	`\.pb\.go$`,
}

type RegExps struct {
	regexps []*regexp.Regexp
	strings []string
}

func MustNewRegExps(exps []string) *RegExps {
	res, _ := NewRegExps(exps)
	if res == nil {
		res = &RegExps{}
	}
	return res
}

func NewRegExps(exps []string) (*RegExps, error) {
	res := &RegExps{
		regexps: make([]*regexp.Regexp, 0, len(exps)),
		strings: make([]string, 0, len(exps)),
	}

	var err error

	for _, exp := range exps {
		e := res.Set(exp)
		if e != nil && err == nil {
			err = e
		}
	}

	return res, err
}

func (p *RegExps) Set(s string) error {
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	p.regexps = append(p.regexps, re)
	p.strings = append(p.strings, s)
	return nil
}

func (p *RegExps) String() string {
	// return strings.Join(p.strings, "\n")
	d, _ := json.Marshal(p.strings)
	return string(d)
	// return fmt.Sprint(p.strings)
}

func (p *RegExps) Check(s string) bool {
	for _, re := range p.regexps {
		loc := re.FindStringIndex(s)
		if len(loc) == 2 {
			return true
		}
	}
	return false
}
