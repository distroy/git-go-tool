/*
 * Copyright (C) distroy
 */

package goformat

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/distroy/git-go-tool/3rd/convey"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/strcore"
)

func Test_labelChecker_Check(t *testing.T) {
	filename := "test.go"

	type args struct {
		data string
	}
	tests := []struct {
		c    importChecker
		args args
		want []*Issue
	}{
		{
			args: args{
				data: strings.Join([]string{
					"package test",
					"",
					"type testStruct struct {",
					"	A int `json:\"a;omitempty\"`",
					"	B int `json: \"b;omitempty\"`",
					"	C int `gorm:\"default:0\"`",
					"}",
				}, "\n"),
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   4,
					EndLine:     4,
					Level:       LevelError,
					Description: `struct field "A" has invalid json label "a;omitempty"`,
				},
				// {
				// 	Filename:    filename,
				// 	BeginLine:   5,
				// 	EndLine:     5,
				// 	Level:       LevelError,
				// 	Description: `invalid "json" label of struct field "B"`,
				// },
				{
					Filename:    filename,
					BeginLine:   6,
					EndLine:     6,
					Level:       LevelError,
					Description: `struct field "C" has invalid column in gorm label`,
				},
			},
		},
		{
			args: args{
				data: strings.Join([]string{
					"package test",
					"",
					"type testStruct struct {",
					"	A int `json:\"a\"`",
					"	B int `json:\"a,omitempty\"`",
					"}",
				}, "\n"),
			},
			want: []*Issue{
				{
					Filename:    filename,
					BeginLine:   3,
					EndLine:     6,
					Level:       LevelError,
					Description: `struct field "B" has duplicate json label "a"`,
				},
			},
		},
		{
			args: args{
				data: strings.Join([]string{
					"package test",
					"",
					"type TestType1 int",
					"type TestType2 string",
					"",
					"type testStruct struct {",
					"	TestType1",
					"	TestType2",
					"}",
				}, "\n"),
			},
			want: []*Issue{
				// {
				// 	Filename:    filename,
				// 	BeginLine:   3,
				// 	EndLine:     6,
				// 	Level:       LevelError,
				// 	Description: `struct field "B" has duplicate json label "a"`,
				// },
			},
		},
	}

	convey.Convey(t.Name(), t, func() {
		for i, tt := range tests {
			name := testGetSubName(t, i)
			// t.Run(tt.name, func(t *testing.T) {
			convey.Convey(name, func() {
				c := LabelChecker(&LabelConfig{
					JsonLabel: true,
				})

				f := filecore.NewTestFile(filename, strcore.StrToBytesUnsafe(tt.args.data))
				x := NewContext(f)

				c.Check(x)

				convey.So(x.Issues(), convey.ShouldResemble, tt.want)
				// if got := x.Issues(); !reflect.DeepEqual(got, tt.want) {
				// 	testPrintCheckResult(t, got, tt.want)
				// }
			})
		}
	})
}

func TestJsonMarshal(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		type People struct {
			PeopleId int `json:"people_id"`
			UserId   int `json:"user_id"`
		}

		type User struct {
			People
			UserId   int    `json:"user_id"`
			UserName string `json:"user_name"`
		}

		user := &User{
			People: People{
				PeopleId: 1,
				UserId:   3,
			},
			UserId:   2,
			UserName: "abc",
		}

		got, err := json.Marshal(user)
		convey.So(err, convey.ShouldBeNil)
		// if err != nil {
		// 	t.Errorf("json marshal error. err:%v", err)
		// }

		// t.Logf("json marshal succ %s", got)

		want := `{"people_id":1,"user_id":2,"user_name":"abc"}`
		convey.So(string(got), convey.ShouldEqual, want)
		// if string(got) != want {
		// 	t.Errorf("json marshal fail. \nwant:%s\ngot:%s", want, got)
		// }
	})
}

func testGetSubName(t testing.TB, idx int) string {
	return fmt.Sprintf("%s/#%02d", t.Name(), idx)
}
