/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/distroy/git-go-tool/core/tagcore"
)

type Flag struct {
	lvl  int
	val  reflect.Value
	tags tagcore.Tags

	Name    string
	Value   Value
	Default string
	Usage   string
	IsArgs  bool
}

type FlagSet struct {
	command *flag.FlagSet
	flags   map[string]*Flag
	args    *Flag
}

func NewFlagSet() *FlagSet {
	s := &FlagSet{
		command: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		flags:   make(map[string]*Flag),
	}

	s.command.Usage = s.printUsage
	return s
}

func (s *FlagSet) printUsage() {
	s.command.VisitAll(func(f *flag.Flag) {
		ff := s.flags[f.Name]
		s.printFlagUsage(ff)
	})
}

func (s *FlagSet) printFlagUsage(f *Flag) {
	const (
		tab           = "    "
		namePrefix    = "  "
		usagePrefix   = "\n  " + tab + tab
		defaultPrefix = usagePrefix + tab + tab
	)

	b := &strings.Builder{}

	fmt.Fprintf(b, "%s-%s", namePrefix, f.Name) // Two spaces before -; see next two comments.
	name, usage := unquoteUsage(f)
	if len(name) > 0 {
		fmt.Fprintf(b, " %s", name)
	}
	// Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if b.Len() <= 4 { // space, space, '-', 'x'.
		fmt.Fprintf(b, tab)
	} else {
		// Four spaces before the tab triggers good alignment
		// for both 4- and 8-space tab stops.
		fmt.Fprint(b, usagePrefix)
	}

	fmt.Fprint(b, strings.ReplaceAll(usage, "\n", usagePrefix))

	if isZeroValue(f.Value, f.Default) {
		fmt.Fprint(s.command.Output(), b.String(), "\n")
		return
	}

	fmt.Fprint(b, usagePrefix, "default: ")
	switch v := f.Value.(type) {
	default:
		if strings.Index(f.Default, "\n") > 0 {
			fmt.Fprint(b, defaultPrefix, strings.ReplaceAll(f.Default, "\n", defaultPrefix))
		} else {
			fmt.Fprintf(b, "%v", f.Default)
		}

	case *stringValue:
		fmt.Fprintf(b, "%q", f.Default)

	case *stringsValue:
		for _, s := range *v {
			fmt.Fprintf(b, "%s%q", defaultPrefix, s)
		}
	}

	fmt.Fprint(s.command.Output(), b.String(), "\n")
}

func (s *FlagSet) MustParse(args ...[]string) {
	a := os.Args[1:]
	if len(args) > 0 {
		a = args[0]
	}

	err := s.parse(a)
	if err != nil {
		panic(fmt.Errorf("parse flag set fail. args:%v, err:%v", a, err))
	}
}

func (s *FlagSet) Parse(args ...[]string) error {
	if len(args) > 0 {
		return s.parse(args[0])
	}

	if err := s.parse(os.Args[1:]); err != nil {
		return err
	}

	if s.args != nil {
		s.args.val.Set(reflect.ValueOf(s.command.Args()))
	}

	return nil
}

func (s *FlagSet) parse(args []string) error {
	err := s.command.Parse(args)
	return err
}

func (s *FlagSet) Model(v interface{}) {
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("input flags must be pointer to struct. %s", typ.String()))
	}
	val = val.Elem()

	s.parseStruct(0, val)
}

func (s *FlagSet) addFlag(f *Flag) {
	if f == nil {
		return
	}

	if f.IsArgs {
		if s.args == nil || s.args.lvl > f.lvl {
			s.args = f
		}
		return
	}

	// if v := s.flags[f.Name]; v != nil {
	// }

	v, val := s.getFlagValue(f)
	if v == nil {
		return
	}
	f.Value = v

	if len(f.Default) > 0 && (val.Kind() != reflect.Slice || val.Len() == 0) {
		v.Set(f.Default)
	}
	f.Default = v.String()

	s.command.Var(v, f.Name, f.Usage)
	s.flags[f.Name] = f
	// log.Printf(" === %s: %v", typ.String(), val.Interface())
}

func (s *FlagSet) getFlagValue(f *Flag) (Value, reflect.Value) {
	val := f.val
	if v, ok := val.Interface().(Value); ok {
		return v, val
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	// log.Printf(" === 2222 %#v", m)

	fn := fillFlagFuncMap[typ]
	if fn == nil {
		return nil, val
	}

	return fn(val), val
}

func (s *FlagSet) parseStruct(lvl int, val reflect.Value) {
	typ := val.Type()

	// log.Printf(" === %s", typ.String())

	for i, l := 0, typ.NumField(); i < l; i++ {
		field := typ.Field(i)
		fVal := val.Field(i)

		if _, ok := fVal.Interface().(Value); ok {
			s.parseFieldFlag(lvl, fVal, field)
			continue
		}

		s.parseStructField(lvl, fVal, field)
	}
}

func (s *FlagSet) parseFieldFlag(lvl int, val reflect.Value, field reflect.StructField) {
	tag, ok := field.Tag.Lookup(tagName)
	if !ok || len(tag) == 0 {
		f := &Flag{
			lvl:  lvl,
			val:  val,
			tags: tagcore.New(),
			Name: parseFlagName(field),
		}
		s.addFlag(f)
		return
	}

	tags := tagcore.Parse(tag)
	if tags.Has("-") {
		return
	}
	// log.Printf(" === %s %#v", tag, tags)

	f := &Flag{
		lvl:     lvl,
		val:     val,
		tags:    tags,
		Name:    tags.Get("name"),
		Usage:   tags.Get("usage"),
		Default: tags.Get("default"),
		IsArgs:  tags.Has("args"),
	}

	if len(f.Name) == 0 {
		f.Name = parseFlagName(field)
	}
	// log.Printf(" === 1111 %#v", m)

	s.addFlag(f)
	return
}

func (s *FlagSet) parseStructField(lvl int, fVal reflect.Value, field reflect.StructField) {
	val := fVal
	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if val.IsNil() {
			val.Set(reflect.New(typ))
		}
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		s.parseStruct(lvl+1, val)
		return
	}

	s.parseFieldFlag(lvl, val, field)
}
