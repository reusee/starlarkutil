package starlarkutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/reusee/e3"
)

type InvalidValue struct {
	Str string
	e3.Prev
}

func (i *InvalidValue) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("invalid value `%s`", i.Str))
	b.WriteString(i.Prev.String("\n"))
	return b.String()
}

func WithInvalidValue(str string) e3.WrapFunc {
	return func(err error) e3.Error {
		return &InvalidValue{
			Str:  str,
			Prev: e3.Prev{Err: err},
		}
	}
}

var TooManyElements = fmt.Errorf("too many elements")

type InvalidTarget struct {
	Str  string
	Type reflect.Type
	e3.Prev
}

func (i *InvalidTarget) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("invalid target `%v` for `%s`", i.Type, i.Str))
	b.WriteString(i.Prev.String("\n"))
	return b.String()
}
