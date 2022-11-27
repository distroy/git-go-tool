package check_test

import (
	. "github.com/distroy/git-go-tool/3rd/yaml/check"
)

var _ = Suite(&PrinterS{})

type PrinterS struct{}

func (s *PrinterS) TestCountSuite(c *C) {
	suitesRun += 1
}

var printTestFuncLine int

func init() {
	printTestFuncLine = getMyLine() + 3
}

func printTestFunc() {
	println(1)  // Comment1
	if 2 == 2 { // Comment2
		println(3) // Comment3
	}
	switch 5 {
	case 6:
		println(7) // Comment7
		println(8)
	}
	switch interface{}(10).(type) { // Comment10
	case int:
		println(12)
		println(13)
	}
	select {
	case <-(chan bool)(nil):
		println(17)
		println(18)
	default:
		println(20)
		println(21)
	}
	println(23,
		24)
	_ = func() {
		println(26)
		println(27)
	}
	println(29, func() {
		println(30)
	})
	// Leading comment
	// with multiple lines.
	println(34) // Comment34
}

var printLineTests = []struct {
	line   int
	output string
}{
	{1, "println(1) // Comment1"},
	{2, "if 2 == 2 { // Comment2\n    ...\n}"},
	{3, "println(3) // Comment3"},
	{5, "switch 5 {\n...\n}"},
	{6, "case 6:\n    ... // Comment7\n    println(8)"},
	{7, "println(7) // Comment7"},
	{8, "println(8)"},
	{10, "switch interface{}(10).(type) { // Comment10\n...\n}"},
	{11, "case int:\n    ...\n    println(13)"},
	{16, "case <-(chan bool)(nil):\n    ...\n    println(18)"},
	{17, "println(17)"},
	{18, "println(18)"},
	{19, "default:\n    ...\n    println(21)"},
	{20, "println(20)"},
	{21, "println(21)"},
	{23, "println(23,\n    24)"},
	{24, "println(23,\n    24)"},
	{25, "_ = func() {\n    println(26)\n    println(27)\n}"},
	{26, "println(26)"},
	{27, "println(27)"},
	{29, "println(29, func() {\n    println(30)\n})"},
	{30, "println(30)"},
	{31, "println(29, func() {\n    println(30)\n})"},
	{34, "// Leading comment\n// with multiple lines.\nprintln(34) // Comment34"},
}

func (s *PrinterS) TestPrintLine(c *C) {
	for _, test := range printLineTests {
		output, err := PrintLine("printer_test.go", printTestFuncLine+test.line)
		c.Assert(err, IsNil)
		c.Assert(output, Equals, test.output)
	}
}

var indentTests = []struct {
	in, out string
}{
	{"", ""},
	{"\n", "\n"},
	{"a", ">>>a"},
	{"a\n", ">>>a\n"},
	{"a\nb", ">>>a\n>>>b"},
	{" ", ">>> "},
}

func (s *PrinterS) TestIndent(c *C) {
	for _, test := range indentTests {
		out := Indent(test.in, ">>>")
		c.Assert(out, Equals, test.out)
	}

}
