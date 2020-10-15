package starlarkutil

import (
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func init() {
	resolve.AllowFloat = true
	resolve.AllowLambda = true
}

func eval(src any) starlark.Value {
	v, err := starlark.Eval(
		new(starlark.Thread),
		"testing",
		src,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return v
}
