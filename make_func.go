package starlarkutil

import (
	"fmt"
	"reflect"
	"strings"

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
	var kwargsSpecs map[string]reflect.StructField
	if numParams > 0 {
		lastType := fnType.In(numParams - 1)
		if lastType.Kind() == reflect.Struct {
			kwargsSpecs = make(map[string]reflect.StructField)
			for i := 0; i < lastType.NumField(); i++ {
				field := lastType.Field(i)
				name := strings.ToLower(field.Name)
				kwargsSpecs[name] = field
			}
		}
	}

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
			if len(kwargsSpecs) > 0 || isVariadic {
				if numArgs < numParams-1 {
					return nil, fmt.Errorf("not enough argument")
				}
			} else {
				return nil, fmt.Errorf("not enough argument")
			}
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
			if i < numArgs {
				// kwargs may omit the last argument
				if err := sb.Copy(
					Marshal(args.Index(i), &paramTypes[i], nil),
					sb.UnmarshalValue(sb.DefaultCtx, ptr, nil),
				); err != nil { // NOCOVER
					return nil, err
				}
			}
			argValues = append(argValues, ptr.Elem())
		}

		if len(kwargs) > 0 {
			kwValue := argValues[len(argValues)-1]
			for _, tuple := range kwargs {
				key := strings.ToLower(tuple[0].(starlark.String).GoString())
				spec, ok := kwargsSpecs[key]
				if !ok {
					return nil, fmt.Errorf("no such keyword argument: %s", key)
				}
				ptr := reflect.New(spec.Type)
				if err := sb.Copy(
					Marshal(tuple[1], &spec.Type, nil),
					sb.UnmarshalValue(sb.DefaultCtx, ptr, nil),
				); err != nil {
					return nil, err
				}
				kwValue.FieldByIndex(spec.Index).Set(ptr.Elem())
			}
		}

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
