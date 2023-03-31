/*
 * Copyright (C) distroy
 */

package goformat

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"github.com/distroy/git-go-tool/core/tagcore"
)

type LabelConfig struct {
	JsonLabel bool
	GormLabel bool
}

func LabelChecker(cfg *LabelConfig) Checker {
	c := labelChecker{
		cfg:            cfg,
		checkFieldList: make([]labelCheckField, 0, 4),
	}

	c.checkFieldList = append(c.checkFieldList, labelCheckField{
		Enable:   cfg.JsonLabel,
		TagName:  "json",
		Function: c.checkJsonField,
	})
	c.checkFieldList = append(c.checkFieldList, labelCheckField{
		Enable:   cfg.GormLabel,
		TagName:  "gorm",
		Function: c.checkGormField,
	})

	return c
}

type labelCheckField struct {
	Enable   bool
	TagName  string
	Function func(x *Context, st *ast.Field, tagName, tabValue string) (string, bool)
}

type labelChecker struct {
	cfg *LabelConfig

	checkFieldList []labelCheckField
}

func (c labelChecker) Check(x *Context) Error {
	file := x.MustParse()

	var err Error
	ast.Inspect(file, func(n ast.Node) bool {
		switch nn := n.(type) {
		case *ast.StructType:
			if nn.Fields == nil {
				return true
			}

			err = c.checkLabels(x, nn)
		}

		return err == nil
	})

	return err
}

func (c labelChecker) checkLabels(x *Context, st *ast.StructType) Error {
	for _, labelType := range c.checkFieldList {
		if !labelType.Enable {
			continue
		}

		err := c.loopCheckFields(x, st, labelType.TagName, labelType.Function)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c labelChecker) loopCheckFields(
	x *Context, st *ast.StructType, tagName string,
	fn func(x *Context, field *ast.Field, tagName, tabValue string) (string, bool),
) Error {

	sPos := x.Position(st.Pos())
	sEnd := x.Position(st.End())

	labels := make(map[string]*ast.Field, len(st.Fields.List))

	for _, field := range st.Fields.List {
		label := ""

		tabValue, ok := c.getStructFieldTag(x, field, tagName)
		if !ok || tabValue == "-" {
			continue

		} else if tabValue == "" {
			label = c.fieldName(x, field)

		} else {
			var ok bool
			label, ok = fn(x, field, tagName, tabValue)
			if !ok {
				continue
			}
		}

		if _, ok := labels[label]; ok {
			// log.Printf("duplicate json tag. [%s] [%s]", c.fieldName(v), c.fieldName(field))
			x.AddIssue(&Issue{
				Filename:  x.Name,
				BeginLine: sPos.Line,
				EndLine:   sEnd.Line,
				Level:     LevelError,
				Description: fmt.Sprintf(`struct field "%s" has duplicate %s label "%s"`,
					c.fieldName(x, field), tagName, label),
			})
			return nil
		}

		labels[label] = field
	}

	return nil
}

func (c labelChecker) getStructFieldTag(x *Context, field *ast.Field, tagName string) (string, bool) {
	// log.Printf("field name. field:%s", c.fieldName(field))

	// 内嵌field
	if len(field.Names) == 0 {
		return "-", true
	}

	if !field.Names[0].IsExported() {
		return "-", true
	}

	if field.Tag == nil {
		return "", true
	}

	// log.Printf("field tag. field:%s, tag:%s", c.fieldName(x, field), field.Tag.Value)

	tag, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		// log.Printf("unable to read struct tag %s", field.Tag.Value)
		fPos := x.Position(field.Pos())
		fEnd := x.Position(field.End())
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: fPos.Line,
			EndLine:   fEnd.Line,
			Level:     LevelError,
			Description: fmt.Sprintf(`unquote "%s" label of struct field "%s" fail`,
				tagName, c.fieldName(x, field)),
		})
		return "", false
	}

	// log.Printf("field tag 2. field:%s, tag:%s", c.fieldName(x, field), tag)

	tagValue, ok := reflect.StructTag(tag).Lookup(tagName)
	if ok && tagValue == "" {
		// fPos := x.Position(field.Pos())
		// fEnd := x.Position(field.End())
		// x.AddIssue(&Issue{
		// 	Filename:  x.Name,
		// 	BeginLine: fPos.Line,
		// 	EndLine:   fEnd.Line,
		// 	Level:     LevelError,
		// 	Description: fmt.Sprintf(`invalid "%s" label of struct field "%s"`,
		// 		tagName, c.fieldName(x, field)),
		// })
		// return "", false
	}

	return tagValue, true
}

func (c labelChecker) fieldName(x *Context, field *ast.Field) string {
	// for i, name := range field.Names {
	// 	log.Printf("field name. index:%d, name:%s", i, name.String())
	// }

	if len(field.Names) > 0 {
		return field.Names[0].String()
	}

	return ""
}

func (c labelChecker) checkJsonField(x *Context, field *ast.Field, tagName, tabValue string) (string, bool) {
	var label string

	arr := strings.Split(tabValue, ",")
	if len(arr) > 0 && arr[0] != "" {
		label = arr[0]
	}

	if label == "" {
		label = c.fieldName(x, field)
	}

	if strings.Contains(label, ";") {
		fPos := x.Position(field.Pos())
		fEnd := x.Position(field.End())
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: fPos.Line,
			EndLine:   fEnd.Line,
			Level:     LevelError,
			Description: fmt.Sprintf(`struct field "%s" has invalid json label "%s"`,
				c.fieldName(x, field), label),
		})
		return "", false
	}

	return label, true
}

func (c labelChecker) checkGormField(x *Context, field *ast.Field, tagName, tagValue string) (string, bool) {
	var label string

	if tagValue == "" {
		label = c.fieldName(x, field)
	}

	tags := tagcore.Parse(tagValue)

	if column := tags.Get("column"); column == "" {
		fPos := x.Position(field.Pos())
		fEnd := x.Position(field.End())
		x.AddIssue(&Issue{
			Filename:  x.Name,
			BeginLine: fPos.Line,
			EndLine:   fEnd.Line,
			Level:     LevelError,
			Description: fmt.Sprintf(`struct field "%s" has invalid column in gorm label`,
				c.fieldName(x, field)),
		})
		return "", false
	}

	return label, true
}
