package starlarkutil

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestMakeFunc(t *testing.T) {
	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i int) int {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i int8) int8 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i int16) int16 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i int32) int32 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i int64) int64 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i uint) uint {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i uint8) uint8 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i uint16) uint16 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i uint32) uint32 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(42) == 84`, starlark.StringDict{
		"foo": MakeFunc("foo", func(i uint64) uint64 {
			return i * 2
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo()`, starlark.StringDict{
		"foo": MakeFunc("foo", func() bool {
			return true
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo((1, 2, 3)) == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(tuple func() (int, int, int)) int {
			a, b, c := tuple()
			return a + b + c
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(set([1, 2, 3])) == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(set map[int]bool) int {
			var sum int
			for e := range set {
				sum += e
			}
			return sum
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo([1, 2, 3]) == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(list []int) int {
			var sum int
			for _, e := range list {
				sum += e
			}
			return sum
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo({'A': 1, 'B': 2}) == 3`, starlark.StringDict{
		"foo": MakeFunc("foo", func(arg struct {
			A int
			B int
		}) int {
			return arg.A + arg.B
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo({'A': 1, 'B': 2}) == 3`, starlark.StringDict{
		"foo": MakeFunc("foo", func(arg map[string]int) int {
			var sum int
			for _, v := range arg {
				sum += v
			}
			return sum
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(3)['N'] == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(l int) (ret struct {
			N int
		}) {
			ret.N = l * 2
			return
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(3)['N'] == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(l int) (ret map[string]int) {
			ret = map[string]int{
				"N": l * 2,
			}
			return
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(3)[0] == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(l int) func() int {
			return func() int {
				return l * 2
			}
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(3)[0] == 6`, starlark.StringDict{
		"foo": MakeFunc("foo", func(l int) []int {
			return []int{l * 2}
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo('foo') == 'foo'`, starlark.StringDict{
		"foo": MakeFunc("foo", func(s string) string {
			return s
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(bytes('foo')) == bytes('foo')`, starlark.StringDict{
		"foo": MakeFunc("foo", func(s []byte) []byte {
			return s
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(2.5) == 2.5`, starlark.StringDict{
		"foo": MakeFunc("foo", func(f float32) float32 {
			return f
		}),
	}).Truth() {
		t.Fatal()
	}

	if !eval(`foo(2.5) == 2.5`, starlark.StringDict{
		"foo": MakeFunc("foo", func(f float64) float64 {
			return f
		}),
	}).Truth() {
		t.Fatal()
	}

}
