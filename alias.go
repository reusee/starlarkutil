package starlarkutil

import (
	"errors"
	"fmt"

	"github.com/reusee/e5"
	"github.com/reusee/sb"
)

type (
	any  = interface{}
	Src  = sb.Proc
	Sink = sb.Sink
)

var (
	is = errors.Is
	as = errors.As
	pt = fmt.Printf
	we = e5.Wrap
	he = e5.Handle
	ce = e5.Check
)
