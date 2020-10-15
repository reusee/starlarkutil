package starlarkutil

import (
	"errors"
	"fmt"

	"github.com/reusee/e3"
)

type (
	any = interface{}
)

var (
	is   = errors.Is
	as   = errors.As
	pt   = fmt.Printf
	wrap = e3.Wrap
)
