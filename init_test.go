package starlarkutil

import "go.starlark.net/resolve"

func init() {
	resolve.AllowSet = true
}
