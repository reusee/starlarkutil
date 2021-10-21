package starlarkutil

import (
	"errors"
	"fmt"

	"github.com/reusee/e4"
)

type (
	any = interface{}
)

var (
	is = errors.Is
	as = errors.As
	pt = fmt.Printf
	we = e4.Wrap
)
