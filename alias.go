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
	is   = errors.Is
	as   = errors.As
	pt   = fmt.Printf
	ce   = e4.Check
	he   = e4.Handle
	wrap = e4.Wrap
)
