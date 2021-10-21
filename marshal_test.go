package starlarkutil

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/reusee/e4"
	"github.com/reusee/sb"
)

func TestMarshal(t *testing.T) {
	defer he(nil, e4.TestingFatal(t))

	var b bool
	ce(sb.Copy(
		Marshal(eval("True"), nil),
		sb.Unmarshal(&b),
	))
	if !b {
		t.Fatal()
	}

	var v any
	ce(sb.Copy(
		Marshal(eval("None"), nil),
		sb.Unmarshal(&v),
	))
	if v != nil {
		t.Fatal()
	}

	var bs []byte
	ce(sb.Copy(
		Marshal(eval("bytes('Foo')"), nil),
		sb.Unmarshal(&bs),
	))
	if !bytes.Equal(bs, []byte("Foo")) {
		t.Fatal()
	}

	var i int
	ce(sb.Copy(
		Marshal(eval("42"), nil),
		sb.Unmarshal(&i),
	))
	if i != 42 {
		t.Fatal()
	}

	var f float64
	ce(sb.Copy(
		Marshal(eval("4.2"), nil),
		sb.Unmarshal(&f),
	))
	if f != 4.2 {
		t.Fatal()
	}

	var s string
	ce(sb.Copy(
		Marshal(eval("'foo'"), nil),
		sb.Unmarshal(&s),
	))
	if s != "foo" {
		t.Fatal()
	}

	var ints []int
	ce(sb.Copy(
		Marshal(eval("[1, 2, 3]"), nil),
		sb.Unmarshal(&ints),
	))
	if len(ints) != 3 {
		t.Fatal()
	}
	if fmt.Sprintf("%v", ints) != "[1 2 3]" {
		t.Fatal()
	}

	var tuple sb.Tuple
	ce(sb.Copy(
		Marshal(eval("(1, 2, 3)"), nil),
		sb.Unmarshal(&tuple),
	))
	if len(ints) != 3 {
		t.Fatal()
	}
	if fmt.Sprintf("%v", tuple) != "[1 2 3]" {
		t.Fatal()
	}

	var m map[int]int
	ce(sb.Copy(
		Marshal(eval("{1: 2, 3: 4, 5: 6}"), nil),
		sb.Unmarshal(&m),
	))
	if len(m) != 3 {
		t.Fatal()
	}
	if m[1] != 2 {
		t.Fatal()
	}
	if m[3] != 4 {
		t.Fatal()
	}
	if m[5] != 6 {
		t.Fatal()
	}

	var set map[int]bool
	ce(sb.Copy(
		Marshal(eval("set((1, 2, 3))"), nil),
		sb.Unmarshal(&set),
	))
	if len(m) != 3 {
		t.Fatal()
	}
	if !set[1] {
		t.Fatal()
	}
	if !set[2] {
		t.Fatal()
	}
	if !set[3] {
		t.Fatal()
	}

}
