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
}
