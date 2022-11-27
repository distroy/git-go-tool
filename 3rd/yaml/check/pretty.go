package check

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/tabwriter"
)

type sbuf []string

func (p *sbuf) Printf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	*p = append(*p, s)
}

// Diff returns a slice where each element describes
// a difference between a and b.
func Diff(a, b interface{}) (desc []string) {
	Pdiff((*sbuf)(&desc), a, b)
	return desc
}

type Printfer interface {
	Printf(format string, a ...interface{})
}

// Pdiff prints to p a description of the differences between a and b.
// It calls Printf once for each difference, with no trailing newline.
// The standard library log.Logger is a Printfer.
func Pdiff(p Printfer, a, b interface{}) {
	diffPrinter{w: p}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
}

type diffPrinter struct {
	w Printfer
	l string // label
}

func (w diffPrinter) printf(f string, a ...interface{}) {
	var l string
	if w.l != "" {
		l = w.l + ": "
	}
	w.w.Printf(l+f, a...)
}

func (w diffPrinter) diff(av, bv reflect.Value) {
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %# v", formatter{v: bv, quote: true})
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%# v != nil", formatter{v: av, quote: true})
		return
	}
	if !av.IsValid() && !bv.IsValid() {
		return
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return
	}

	switch kind := at.Kind(); kind {
	case reflect.Bool:
		if a, b := av.Bool(), bv.Bool(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if a, b := av.Int(), bv.Int(); a != b {
			w.printf("%d != %d", a, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if a, b := av.Uint(), bv.Uint(); a != b {
			w.printf("%d != %d", a, b)
		}
	case reflect.Float32, reflect.Float64:
		if a, b := av.Float(), bv.Float(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Complex64, reflect.Complex128:
		if a, b := av.Complex(), bv.Complex(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Array:
		n := av.Len()
		for i := 0; i < n; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		if a, b := av.Pointer(), bv.Pointer(); a != b {
			w.printf("%#x != %#x", a, b)
		}
	case reflect.Interface:
		w.diff(av.Elem(), bv.Elem())
	case reflect.Map:
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.printf("%q != (missing)", av.MapIndex(k))
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.diff(av.MapIndex(k), bv.MapIndex(k))
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.printf("(missing) != %q", bv.MapIndex(k))
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %# v", formatter{v: bv, quote: true})
		case !av.IsNil() && bv.IsNil():
			w.printf("%# v != nil", formatter{v: av, quote: true})
		case !av.IsNil() && !bv.IsNil():
			w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Slice:
		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
			break
		}
		for i := 0; i < lenA; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
	case reflect.String:
		if a, b := av.String(), bv.String(); a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i))
		}
	default:
		panic("unknown reflect Kind: " + kind.String())
	}
}

func (d diffPrinter) relabel(name string) (d1 diffPrinter) {
	d1 = d
	if d.l != "" && name[0] != '[' {
		d1.l += "."
	}
	d1.l += name
	return d1
}

// keyEqual compares a and b for equality.
// Both a and b must be valid map keys.
func keyEqual(av, bv reflect.Value) bool {
	if !av.IsValid() && !bv.IsValid() {
		return true
	}
	if !av.IsValid() || !bv.IsValid() || av.Type() != bv.Type() {
		return false
	}
	switch kind := av.Kind(); kind {
	case reflect.Bool:
		a, b := av.Bool(), bv.Bool()
		return a == b
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		a, b := av.Int(), bv.Int()
		return a == b
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		a, b := av.Uint(), bv.Uint()
		return a == b
	case reflect.Float32, reflect.Float64:
		a, b := av.Float(), bv.Float()
		return a == b
	case reflect.Complex64, reflect.Complex128:
		a, b := av.Complex(), bv.Complex()
		return a == b
	case reflect.Array:
		for i := 0; i < av.Len(); i++ {
			if !keyEqual(av.Index(i), bv.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.UnsafePointer, reflect.Ptr:
		a, b := av.Pointer(), bv.Pointer()
		return a == b
	case reflect.Interface:
		return keyEqual(av.Elem(), bv.Elem())
	case reflect.String:
		a, b := av.String(), bv.String()
		return a == b
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			if !keyEqual(av.Field(i), bv.Field(i)) {
				return false
			}
		}
		return true
	default:
		panic("invalid map key type " + av.Type().String())
	}
}

func keyDiff(a, b []reflect.Value) (ak, both, bk []reflect.Value) {
	for _, av := range a {
		inBoth := false
		for _, bv := range b {
			if keyEqual(av, bv) {
				inBoth = true
				both = append(both, av)
				break
			}
		}
		if !inBoth {
			ak = append(ak, av)
		}
	}
	for _, bv := range b {
		inBoth := false
		for _, av := range a {
			if keyEqual(av, bv) {
				inBoth = true
				break
			}
		}
		if !inBoth {
			bk = append(bk, bv)
		}
	}
	return
}

type formatter struct {
	v     reflect.Value
	force bool
	quote bool
}

// printValue must keep track of already-printed pointer values to avoid
// infinite recursion.
type visit struct {
	v   uintptr
	typ reflect.Type
}

// Formatter makes a wrapper, f, that will format x as go source with line
// breaks and tabs. Object f responds to the "%v" formatting verb when both the
// "#" and " " (space) flags are set, for example:
//
//	fmt.Sprintf("%# v", Formatter(x))
//
// If one of these two flags is not set, or any other verb is used, f will
// format x according to the usual rules of package fmt.
// In particular, if x satisfies fmt.Formatter, then x.Format will be called.
func Formatter(x interface{}) (f fmt.Formatter) {
	return formatter{v: reflect.ValueOf(x), quote: true}
}

func (fo formatter) String() string {
	return fmt.Sprint(fo.v.Interface()) // unwrap it
}

func (fo formatter) passThrough(f fmt.State, c rune) {
	s := "%"
	for i := 0; i < 128; i++ {
		if f.Flag(i) {
			s += string(rune(i))
		}
	}
	if w, ok := f.Width(); ok {
		s += fmt.Sprintf("%d", w)
	}
	if p, ok := f.Precision(); ok {
		s += fmt.Sprintf(".%d", p)
	}
	s += string(c)
	fmt.Fprintf(f, s, fo.v.Interface())
}

func (fo formatter) Format(f fmt.State, c rune) {
	if fo.force || c == 'v' && f.Flag('#') && f.Flag(' ') {
		w := tabwriter.NewWriter(f, 4, 4, 1, ' ', 0)
		p := &prettryPrinter{tw: w, Writer: w, visited: make(map[visit]int)}
		p.printValue(fo.v, true, fo.quote)
		w.Flush()
		return
	}
	fo.passThrough(f, c)
}

type prettryPrinter struct {
	io.Writer
	tw      *tabwriter.Writer
	visited map[visit]int
	depth   int
}

func (p *prettryPrinter) indent() *prettryPrinter {
	q := *p
	q.tw = tabwriter.NewWriter(p.Writer, 4, 4, 1, ' ', 0)
	q.Writer = NewIndentWriter(q.tw, []byte{'\t'})
	return &q
}

func (p *prettryPrinter) printInline(v reflect.Value, x interface{}, showType bool) {
	if showType {
		io.WriteString(p, v.Type().String())
		fmt.Fprintf(p, "(%#v)", x)
	} else {
		fmt.Fprintf(p, "%#v", x)
	}
}

func (p *prettryPrinter) printValue(v reflect.Value, showType, quote bool) {
	if p.depth > 10 {
		io.WriteString(p, "!%v(DEPTH EXCEEDED)")
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		p.printInline(v, v.Bool(), showType)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.printInline(v, v.Int(), showType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.printInline(v, v.Uint(), showType)
	case reflect.Float32, reflect.Float64:
		p.printInline(v, v.Float(), showType)
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(p, "%#v", v.Complex())
	case reflect.String:
		p.fmtString(v.String(), quote)
	case reflect.Map:
		t := v.Type()
		if showType {
			io.WriteString(p, t.String())
		}
		writeByte(p, '{')
		if nonzero(v) {
			expand := !canInline(v.Type())
			pp := p
			if expand {
				writeByte(p, '\n')
				pp = p.indent()
			}
			keys := v.MapKeys()
			for i := 0; i < v.Len(); i++ {
				k := keys[i]
				mv := v.MapIndex(k)
				pp.printValue(k, false, true)
				writeByte(pp, ':')
				if expand {
					writeByte(pp, '\t')
				}
				showTypeInStruct := t.Elem().Kind() == reflect.Interface
				pp.printValue(mv, showTypeInStruct, true)
				if expand {
					io.WriteString(pp, ",\n")
				} else if i < v.Len()-1 {
					io.WriteString(pp, ", ")
				}
			}
			if expand {
				pp.tw.Flush()
			}
		}
		writeByte(p, '}')
	case reflect.Struct:
		t := v.Type()
		if v.CanAddr() {
			addr := v.UnsafeAddr()
			vis := visit{addr, t}
			if vd, ok := p.visited[vis]; ok && vd < p.depth {
				p.fmtString(t.String()+"{(CYCLIC REFERENCE)}", false)
				break // don't print v again
			}
			p.visited[vis] = p.depth
		}

		if showType {
			io.WriteString(p, t.String())
		}
		writeByte(p, '{')
		if nonzero(v) {
			expand := !canInline(v.Type())
			pp := p
			if expand {
				writeByte(p, '\n')
				pp = p.indent()
			}
			for i := 0; i < v.NumField(); i++ {
				showTypeInStruct := true
				if f := t.Field(i); f.Name != "" {
					io.WriteString(pp, f.Name)
					writeByte(pp, ':')
					if expand {
						writeByte(pp, '\t')
					}
					showTypeInStruct = labelType(f.Type)
				}
				pp.printValue(getField(v, i), showTypeInStruct, true)
				if expand {
					io.WriteString(pp, ",\n")
				} else if i < v.NumField()-1 {
					io.WriteString(pp, ", ")
				}
			}
			if expand {
				pp.tw.Flush()
			}
		}
		writeByte(p, '}')
	case reflect.Interface:
		switch e := v.Elem(); {
		case e.Kind() == reflect.Invalid:
			io.WriteString(p, "nil")
		case e.IsValid():
			pp := *p
			pp.depth++
			pp.printValue(e, showType, true)
		default:
			io.WriteString(p, v.Type().String())
			io.WriteString(p, "(nil)")
		}
	case reflect.Array, reflect.Slice:
		t := v.Type()
		if showType {
			io.WriteString(p, t.String())
		}
		if v.Kind() == reflect.Slice && v.IsNil() && showType {
			io.WriteString(p, "(nil)")
			break
		}
		if v.Kind() == reflect.Slice && v.IsNil() {
			io.WriteString(p, "nil")
			break
		}
		writeByte(p, '{')
		expand := !canInline(v.Type())
		pp := p
		if expand {
			writeByte(p, '\n')
			pp = p.indent()
		}
		for i := 0; i < v.Len(); i++ {
			showTypeInSlice := t.Elem().Kind() == reflect.Interface
			pp.printValue(v.Index(i), showTypeInSlice, true)
			if expand {
				io.WriteString(pp, ",\n")
			} else if i < v.Len()-1 {
				io.WriteString(pp, ", ")
			}
		}
		if expand {
			pp.tw.Flush()
		}
		writeByte(p, '}')
	case reflect.Ptr:
		e := v.Elem()
		if !e.IsValid() {
			writeByte(p, '(')
			io.WriteString(p, v.Type().String())
			io.WriteString(p, ")(nil)")
		} else {
			pp := *p
			pp.depth++
			writeByte(pp, '&')
			pp.printValue(e, true, true)
		}
	case reflect.Chan:
		x := v.Pointer()
		if showType {
			writeByte(p, '(')
			io.WriteString(p, v.Type().String())
			fmt.Fprintf(p, ")(%#v)", x)
		} else {
			fmt.Fprintf(p, "%#v", x)
		}
	case reflect.Func:
		io.WriteString(p, v.Type().String())
		io.WriteString(p, " {...}")
	case reflect.UnsafePointer:
		p.printInline(v, v.Pointer(), showType)
	case reflect.Invalid:
		io.WriteString(p, "nil")
	}
}

func canInline(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map:
		return !canExpand(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if canExpand(t.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Interface:
		return false
	case reflect.Array, reflect.Slice:
		return !canExpand(t.Elem())
	case reflect.Ptr:
		return false
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return false
	}
	return true
}

func canExpand(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map, reflect.Struct,
		reflect.Interface, reflect.Array, reflect.Slice,
		reflect.Ptr:
		return true
	}
	return false
}

func labelType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface, reflect.Struct:
		return true
	}
	return false
}

func (p *prettryPrinter) fmtString(s string, quote bool) {
	if quote {
		s = strconv.Quote(s)
	}
	io.WriteString(p, s)
}

func writeByte(w io.Writer, b byte) {
	w.Write([]byte{b})
}

func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		val = val.Elem()
	}
	return val
}

func nonzero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() != complex(0, 0)
	case reflect.String:
		return v.String() != ""
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if nonzero(getField(v, i)) {
				return true
			}
		}
		return false
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if nonzero(v.Index(i)) {
				return true
			}
		}
		return false
	case reflect.Map, reflect.Interface, reflect.Slice, reflect.Ptr, reflect.Chan, reflect.Func:
		return !v.IsNil()
	case reflect.UnsafePointer:
		return v.Pointer() != 0
	}
	return true
}

// Writer indents each line of its input.
type indentWriter struct {
	w   io.Writer
	bol bool
	pre [][]byte
	sel int
	off int
}

// NewIndentWriter makes a new write filter that indents the input
// lines. Each line is prefixed in order with the corresponding
// element of pre. If there are more lines than elements, the last
// element of pre is repeated for each subsequent line.
func NewIndentWriter(w io.Writer, pre ...[]byte) io.Writer {
	return &indentWriter{
		w:   w,
		pre: pre,
		bol: true,
	}
}

// The only errors returned are from the underlying indentWriter.
func (w *indentWriter) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if w.bol {
			var i int
			i, err = w.w.Write(w.pre[w.sel][w.off:])
			w.off += i
			if err != nil {
				return n, err
			}
		}
		_, err = w.w.Write([]byte{c})
		if err != nil {
			return n, err
		}
		n++
		w.bol = c == '\n'
		if w.bol {
			w.off = 0
			if w.sel < len(w.pre)-1 {
				w.sel++
			}
		}
	}
	return n, nil
}
