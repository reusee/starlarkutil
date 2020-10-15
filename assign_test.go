package starlarkutil

import (
	"reflect"
	"testing"
)

func TestAssign(t *testing.T) {

	// bool
	var b bool
	if err := Assign(eval("True"), reflect.ValueOf(&b)); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal()
	}

	// int
	var i8 int8
	if err := Assign(eval("42"), reflect.ValueOf(&i8)); err != nil {
		t.Fatal(err)
	}
	if i8 != 42 {
		t.Fatal()
	}

	// uint
	var u16 uint16
	if err := Assign(eval("42"), reflect.ValueOf(&u16)); err != nil {
		t.Fatal(err)
	}
	if u16 != 42 {
		t.Fatal()
	}

	// float
	var f32 float32
	if err := Assign(eval("42"), reflect.ValueOf(&f32)); err != nil {
		t.Fatal(err)
	}
	if f32 != 42 {
		t.Fatal()
	}

	// array
	var array [2]int
	if err := Assign(eval("[42, 1]"), reflect.ValueOf(&array)); err != nil {
		t.Fatal(err)
	}
	if array[0] != 42 {
		t.Fatal()
	}
	if array[1] != 1 {
		t.Fatal()
	}

	err := Assign(eval("[42, 1, 2]"), reflect.ValueOf(&array))
	if !is(err, TooManyElements) {
		t.Fatal()
	}
	var errInvalidValue *InvalidValue
	if !as(err, &errInvalidValue) {
		t.Fatal()
	}

	// interface
	var v any
	if err := Assign(eval("42.2"), reflect.ValueOf(&v)); err != nil {
		t.Fatal()
	}
	if v != 42.2 {
		t.Fatal()
	}
	if err := Assign(eval("42"), reflect.ValueOf(&v)); err != nil {
		t.Fatal()
	}
	if v != 42 {
		t.Fatal()
	}
	if err := Assign(eval(`"42"`), reflect.ValueOf(&v)); err != nil {
		t.Fatal()
	}
	if v != "42" {
		t.Fatal()
	}
	if err := Assign(eval(`[1, 2, 3]`), reflect.ValueOf(&v)); err != nil {
		t.Fatal(err)
	}
	if l, ok := v.([]any); !ok {
		t.Fatal()
	} else if len(l) != 3 {
		t.Fatal()
	} else if l[2] != 3 {
		t.Fatal()
	}
	if err := Assign(eval(`{1: 2, 3: 4}`), reflect.ValueOf(&v)); err != nil {
		t.Fatal()
	}
	if m, ok := v.(map[any]any); !ok {
		t.Fatal()
	} else if m[1] != 2 {
		t.Fatal()
	}
	if err := Assign(eval(`True`), reflect.ValueOf(&v)); err != nil {
		t.Fatal(err)
	}
	if v != true {
		t.Fatal()
	}

	// map
	var m map[int]string
	if err := Assign(eval(`{1: "foo", 3: "bar"}`), reflect.ValueOf(&m)); err != nil {
		t.Fatal()
	}
	if m[1] != "foo" {
		t.Fatal()
	}
	if m[3] != "bar" {
		t.Fatal()
	}

	// pointer
	var ip ***int
	if err := Assign(eval(`42`), reflect.ValueOf(&ip)); err != nil {
		t.Fatal(err)
	}
	if ***ip != 42 {
		t.Fatal()
	}
	var p **int
	ip = &p
	if err := Assign(eval(`1`), reflect.ValueOf(&ip)); err != nil {
		t.Fatal(err)
	}
	if ***ip != 1 {
		t.Fatal()
	}

	// slice
	var slice []int
	if err := Assign(eval(`[1, 2, 3, 4]`), reflect.ValueOf(&slice)); err != nil {
		t.Fatal(err)
	}
	if len(slice) != 4 {
		t.Fatal()
	}
	if slice[3] != 4 {
		t.Fatal()
	}
	var anySlice []any
	if err := Assign(eval(`[1, True, "foo"]`), reflect.ValueOf(&anySlice)); err != nil {
		t.Fatal(err)
	}
	if len(anySlice) != 3 {
		t.Fatal()
	}
	if anySlice[2] != "foo" {
		t.Fatal()
	}

	// string
	var s string
	if err := Assign(eval(`"foo"`), reflect.ValueOf(&s)); err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal()
	}

	// struct
	var st struct {
		Foo int
		Bar string
	}
	if err := Assign(eval(`{
	  "Foo": 42,
	  "Bar": "bar",
	  "Baz": 42.0,
    42: 42,
	}`), reflect.ValueOf(&st)); err != nil {
		t.Fatal(err)
	}
	if st.Foo != 42 {
		t.Fatal()
	}
	if st.Bar != "bar" {
		t.Fatal()
	}

	// nil target
	if err := Assign(eval("42"), reflect.ValueOf((*int)(nil))); err != nil {
		t.Fatal(err)
	}

}

func TestAssignInvalidTarget(t *testing.T) {

	// bad target
	err := Assign(eval("42"), reflect.ValueOf(0))
	var invalidTarget *InvalidTarget
	if !as(err, &invalidTarget) {
		t.Fatal()
	}

	// type mismatch
	var ip ***int
	err = Assign(eval("True"), reflect.ValueOf(&ip))
	var invalidValue *InvalidValue
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// ivnalid type
	var c chan bool
	err = Assign(eval("42"), reflect.ValueOf(&c))
	if !as(err, &invalidTarget) {
		t.Fatal()
	}

	// invalid int
	var i int
	err = Assign(eval("True"), reflect.ValueOf(&i))
	if !as(err, &invalidValue) {
		t.Fatal()
	}
	err = Assign(eval("100000000000000000000000"), reflect.ValueOf(&i))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid uint
	var u uint
	err = Assign(eval("True"), reflect.ValueOf(&u))
	if !as(err, &invalidValue) {
		t.Fatal()
	}
	err = Assign(eval("100000000000000000000000"), reflect.ValueOf(&u))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid float
	var f float64
	err = Assign(eval("True"), reflect.ValueOf(&f))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid array
	var ia [2]int
	err = Assign(eval("[1, True]"), reflect.ValueOf(&ia))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid interface
	var iface interface {
		Foo()
	}
	err = Assign(eval("lambda: 42"), reflect.ValueOf(&iface))
	if !as(err, &invalidValue) {
		t.Fatal()
	}
	var some any
	err = Assign(eval("10000000000000000000000000"), reflect.ValueOf(&some))
	if !as(err, &invalidValue) {
		t.Fatal()
	}
	err = Assign(eval("[10000000000000000000000000]"), reflect.ValueOf(&some))
	if !as(err, &invalidValue) {
		t.Fatal()
	}
	err = Assign(eval("{1: 10000000000000000000000000}"), reflect.ValueOf(&some))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid map
	var m map[int]string
	err = Assign(eval("{10000000000000000000000000: 1}"), reflect.ValueOf(&m))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid string
	var s string
	err = Assign(eval("42"), reflect.ValueOf(&s))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

	// invalid struct
	var st struct {
		Foo int
	}
	err = Assign(eval("{'Foo': 10000000000000000000000000}"), reflect.ValueOf(&st))
	if !as(err, &invalidValue) {
		t.Fatal()
	}

}
