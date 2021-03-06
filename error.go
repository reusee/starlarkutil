package starlarkutil

import (
	"fmt"
	"reflect"

	"github.com/reusee/e4"
)

type InvalidValue struct {
	Str string
}

func (i *InvalidValue) Error() string {
	return fmt.Sprintf("invalid value `%s`", i.Str)
}

func WithInvalidValue(str string) e4.WrapFunc {
	return func(err error) error {
		return e4.Error{
			Err: &InvalidValue{
				Str: str,
			},
			Prev: err,
		}
	}
}

var TooManyElements = fmt.Errorf("too many elements")

type InvalidTarget struct {
	Type reflect.Type
	Str  string
}

func (i *InvalidTarget) Error() string {
	return fmt.Sprintf("invalid target `%v` for `%s`", i.Type, i.Str)
}
