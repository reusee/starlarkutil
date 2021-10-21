package starlarkutil

import (
	"errors"
	"fmt"

	"github.com/reusee/e4"
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
	we = e4.Wrap
	he = e4.Handle
	ce = e4.Check
)
