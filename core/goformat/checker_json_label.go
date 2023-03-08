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
)

func JsonLabelChecker(enable bool) Checker {
	if !enable {
		return checkerNil{}
	}
	return jsonLabelChecker{}
}

type jsonLabelChecker struct {
	// A int `json:"a"`
	// B int `json:"a"`
}

func (c jsonLabelChecker) Check(x *Context) Error {
	file := x.MustParse()

	var err Error
	ast.Inspect(file, func(n ast.Node) bool {
		switch nn := n.(type) {
		case *ast.StructType:
			err = c.checkStructJsonLabel(x, nn)
		}

		return err == nil
	})

	return err
}

func (c jsonLabelChecker) checkStructJsonLabel(x *Context, st *ast.StructType) Error {
	if st.Fields == nil {
		return nil
	}

	sPos := x.Position(st.Pos())
	sEnd := x.Position(st.End())

	labels := make(map[string]*ast.Field, len(st.Fields.List))

	for _, field := range st.Fields.List {
		label, ok := c.parseStructFieldTagName(x, field)
		if !ok {
			continue

		} else if label == "" {
			break

		} else if label == "-" {
			continue
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
			return nil
		}

		if _, ok := labels[label]; ok {
			// log.Printf("duplicate json tag. [%s] [%s]", c.fieldName(v), c.fieldName(field))
			x.AddIssue(&Issue{
				Filename:  x.Name,
				BeginLine: sPos.Line,
				EndLine:   sEnd.Line,
				Level:     LevelError,
				Description: fmt.Sprintf(`struct field "%s" has duplicate json label "%s"`,
					c.fieldName(x, field), label),
			})
			return nil
		}

		labels[label] = field
	}

	return nil
}

func (c jsonLabelChecker) parseStructFieldTagName(x *Context, field *ast.Field) (string, bool) {
	// log.Printf("field name. field:%s", c.fieldName(field))

	// 内嵌field
	if len(field.Names) == 0 {
		return "-", true
	}

	if !field.Names[0].IsExported() {
		return "-", true
	}

	if field.Tag == nil {
		return c.fieldName(x, field), true
	}

	// log.Printf("field tag. field:%s, tag:%s", c.fieldName(field), field.Tag.Value)

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
			Description: fmt.Sprintf(`unquote json label of struct field "%s" fail`,
				c.fieldName(x, field)),
		})
		return "", false
	}

	// log.Printf("field tag 2. field:%s, tag:%s", c.fieldName(field), tag)

	jsonTag := reflect.StructTag(tag).Get("json")
	if jsonTag == "" {
		return c.fieldName(x, field), true
	}

	// log.Printf("field json tag. field:%s, json tag:%s", c.fieldName(field), jsonTag)

	arr := strings.Split(jsonTag, ",")
	if len(arr) > 0 && arr[0] != "" {
		return arr[0], true
	}

	return c.fieldName(x, field), true
}

func (c jsonLabelChecker) fieldName(x *Context, field *ast.Field) string {
	// for i, name := range field.Names {
	// 	log.Printf("field name. index:%d, name:%s", i, name.String())
	// }

	if len(field.Names) > 0 {
		return field.Names[0].String()
	}

	return ""
}
