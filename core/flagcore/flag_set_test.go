/*
 * Copyright (C) distroy
 */

package flagcore

import (
	"reflect"
	"strings"
	"testing"

	"github.com/distroy/git-go-tool/3rd/convey"
	"github.com/distroy/git-go-tool/core/ptrcore"
)

func TestNewFlagSet(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("", func() {
			s := NewFlagSet()
			convey.So(s, convey.ShouldNotBeNil)
		})
	})
}

func TestFlagSet_printUsageHeader(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey(`name == "" && no args`, func() {
			s := &FlagSet{
				name: "",
				args: nil,
			}
			b := &strings.Builder{}
			s.printUsageHeader(b)
			convey.So(b.String(), convey.ShouldEqual, "Usage of <command>:\nFlags:\n")
		})
		convey.Convey(`name == "abc" && args.meta == ""`, func() {
			s := &FlagSet{
				name: "abc",
				args: &Flag{Meta: ""},
			}
			b := &strings.Builder{}
			s.printUsageHeader(b)
			convey.So(b.String(), convey.ShouldEqual, "Usage: abc [<flags>] [<arg>...]\n\nFlags:\n")
		})
	})
}

func TestFlagSet_Model(t *testing.T) {
	testClearFlagSet := func(s *FlagSet) {
		for _, v := range s.flagSlice {
			v.tags = nil
			v.val = reflect.Value{}
			v.Value = nil
		}
		if args := s.args; args != nil {
			args.tags = nil
			args.val = reflect.Value{}
			args.Value = nil
		}
	}

	type Flags struct {
		Top      int      `flag:"name:top; meta:N; usage:show the top <N>"`
		Avg      bool     `flag:"usage:show the average complexity"`
		DebugLog bool     `flag:"usage:print debug log; bool"`
		Rate     float64  `flag:"default:0.65; usage:"`
		Branch   string   `flag:"meta:branch; usage:git branch name"`
		Pathes   []string `flag:"args; meta:path; default:."`
	}

	convey.Convey(t.Name(), t, func() {
		convey.Convey("", func() {
			flags := &Flags{}

			s := NewFlagSet()
			s.Model(flags)
			testClearFlagSet(s)

			convey.So(s.flagSlice, convey.ShouldResemble, []*Flag{
				{
					lvl:     0,
					Name:    "top",
					Meta:    "N",
					Default: "0",
					Usage:   "show the top <N>",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "avg",
					Meta:    "",
					Default: "false",
					Usage:   "show the average complexity",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "debug-log",
					Meta:    "",
					Default: "false",
					Usage:   "print debug log",
					IsArgs:  false,
					Bool:    true,
				},
				{
					lvl:     0,
					Name:    "rate",
					Meta:    "",
					Default: "0.65",
					Usage:   "",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "branch",
					Meta:    "branch",
					Default: "",
					Usage:   "git branch name",
					IsArgs:  false,
					Bool:    false,
				},
			})
			convey.So(s.args, convey.ShouldResemble, &Flag{
				lvl:     0,
				Name:    "pathes",
				Meta:    "path",
				Default: ".",
				Usage:   "",
				IsArgs:  true,
				Bool:    false,
			})
		})
	})
}

func TestFlagSet_Parse(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.Convey("normal value", func() {
			type Flags struct {
				Top      int      `flag:"name:top; meta:N; usage:show the top <N>"`
				Avg      bool     `flag:"usage:show the average complexity"`
				DebugLog bool     `flag:"usage:print debug log; bool"`
				Rate     float64  `flag:"default:0.65; usage:"`
				Branch   string   `flag:"meta:branch; usage:git branch name"`
				Pathes   []string `flag:"args; meta:path; default:."`
			}

			convey.Convey("no set default", func() {
				flags := &Flags{}
				s := NewFlagSet()
				s.noDefault = true
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				convey.So(err, convey.ShouldBeNil)
				convey.So(flags, convey.ShouldResemble, &Flags{
					Top:      5,
					Avg:      true,
					DebugLog: true,
					Rate:     0,
					Branch:   "",
					Pathes:   []string{},
				})
			})

			convey.Convey("set default", func() {
				flags := &Flags{}
				s := NewFlagSet()
				s.noDefault = false
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				convey.So(err, convey.ShouldBeNil)
				convey.So(flags, convey.ShouldResemble, &Flags{
					Top:      5,
					Avg:      true,
					DebugLog: true,
					Rate:     0.65,
					Branch:   "",
					Pathes:   []string{"."},
				})
			})
		})

		convey.Convey("ptr value", func() {
			type Flags struct {
				Over     *int     `flag:"name:over; meta:N; default:15; usage:show functions with complexity <N>"`
				Top      *int     `flag:"name:top; meta:N; usage:show the top <N>"`
				Avg      *bool    `flag:"usage:show the average complexity"`
				DebugLog *bool    `flag:"usage:print debug log; bool"`
				Rate     *float64 `flag:"default:0.65; usage:"`
				Branch   *string  `flag:"meta:branch; usage:git branch name"`
				Pathes   []string `flag:"args; meta:path; default:."`
			}

			convey.Convey("no set default", func() {
				flags := &Flags{}
				s := NewFlagSet()
				s.noDefault = true
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				convey.So(err, convey.ShouldBeNil)
				convey.So(flags, convey.ShouldResemble, &Flags{
					Over:     nil,
					Top:      ptrcore.NewInt(5),
					Avg:      ptrcore.NewBool(true),
					DebugLog: ptrcore.NewBool(true),
					Rate:     nil,
					Branch:   nil,
					Pathes:   []string{},
				})
			})

			convey.Convey("set default", func() {
				flags := &Flags{}
				s := NewFlagSet()
				s.noDefault = false
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				convey.So(err, convey.ShouldBeNil)
				convey.So(flags, convey.ShouldResemble, &Flags{
					Over:     ptrcore.NewInt(15),
					Top:      ptrcore.NewInt(5),
					Avg:      ptrcore.NewBool(true),
					DebugLog: ptrcore.NewBool(true),
					Rate:     ptrcore.NewFloat64(0.65),
					Branch:   nil,
					Pathes:   []string{"."},
				})
			})
		})
	})
}
