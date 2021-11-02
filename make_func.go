package starlarkutil

import (
	"fmt"
	"reflect"

	"github.com/reusee/sb"
	"go.starlark.net/starlark"
)

func MakeFunc(name string, fn any) *starlark.Builtin {

	//TODO cache info
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func { // NOCOVER
		panic(fmt.Errorf("not a function: %T", fn))
	}
	numParams := fnType.NumIn()
	numReturn := fnType.NumOut()
	var paramTypes []reflect.Type
	for i := 0; i < numParams; i++ {
		paramTypes = append(paramTypes, fnType.In(i))
	}
	isVariadic := fnType.IsVariadic()

	return starlark.NewBuiltin(name, func(
		t *starlark.Thread,
		builtinFunc *starlark.Builtin,
		args starlark.Tuple,
		kwargs []starlark.Tuple,
	) (
		ret starlark.Value,
		err error,
	) {

		numArgs := args.Len()
		if numArgs < numParams { // NOCOVER
			return nil, fmt.Errorf("not enough argument")
		}

		var argValues []reflect.Value
		for i := 0; i < numParams; i++ {

			if isVariadic && i == numParams-1 {
				t := paramTypes[i].Elem()
				for ; i < numArgs; i++ {
					ptr := reflect.New(t)
					if err := sb.Copy(
						Marshal(args.Index(i), &t, nil),
						sb.UnmarshalValue(sb.DefaultCtx, ptr, nil),
					); err != nil { // NOCOVER
						return nil, err
					}
					argValues = append(argValues, ptr.Elem())
				}
				break
			}

			ptr := reflect.New(paramTypes[i])
			if err := sb.Copy(
				Marshal(args.Index(i), &paramTypes[i], nil),
				sb.UnmarshalValue(sb.DefaultCtx, ptr, nil),
			); err != nil { // NOCOVER
				return nil, err
			}
			argValues = append(argValues, ptr.Elem())
		}

		//TODO kwargs

		retValues := fnValue.Call(argValues)

		if numReturn == 0 {
			ret = starlark.None
			return
		}

		if numReturn == 1 {
			proc := sb.MarshalValue(sb.DefaultCtx, retValues[0], nil)
			if err := sb.Copy(
				&proc,
				Unmarshal(&ret, nil),
			); err != nil { // NOCOVER
				return nil, err
			}
			return
		}

		var tuple starlark.Tuple
		for _, value := range retValues {
			var elem starlark.Value
			proc := sb.MarshalValue(sb.DefaultCtx, value, nil)
			if err := sb.Copy(
				&proc,
				Unmarshal(&elem, nil),
			); err != nil { // NOCOVER
				return nil, err
			}
			tuple = append(tuple, elem)
		}
		ret = tuple

		return
	})
}
