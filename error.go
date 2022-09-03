package starlarkutil

import (
	"fmt"
	"reflect"

	"github.com/reusee/e5"
)

type InvalidValue struct {
	Str string
}

func (i *InvalidValue) Error() string {
	return fmt.Sprintf("invalid value `%s`", i.Str)
}

func WithInvalidValue(str string) e5.WrapFunc {
	return func(err error) error {
		return e5.Chain(
			&InvalidValue{
				Str: str,
			},
			err,
		)
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
