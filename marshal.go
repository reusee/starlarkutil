package starlarkutil

import (
	"fmt"
	"reflect"

	"github.com/reusee/sb"
	"go.starlark.net/starlark"
)

func Marshal(value starlark.Value, t *reflect.Type, cont Src) *Src {
	src := MarshalValue(value, t, cont)
	return &src
}

func MarshalValue(value starlark.Value, t *reflect.Type, cont Src) Src {
	return func() (*sb.Token, Src, error) {

		switch value := value.(type) {

		case starlark.NoneType:
			return &sb.Token{
				Kind: sb.KindNil,
			}, cont, nil

		case starlark.Bool:
			return &sb.Token{
				Kind:  sb.KindBool,
				Value: bool(value),
			}, cont, nil

		case starlark.Bytes:
			return &sb.Token{
				Kind:  sb.KindBytes,
				Value: []byte(value),
			}, cont, nil

		case starlark.Int:
			if t == nil {
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt,
					Value: int(i64),
				}, cont, nil
			}

			switch (*t).Kind() {
			case reflect.Int:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt,
					Value: int(i64),
				}, cont, nil
			case reflect.Int8:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt8,
					Value: int8(i64),
				}, cont, nil
			case reflect.Int16:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt16,
					Value: int16(i64),
				}, cont, nil
			case reflect.Int32:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt32,
					Value: int32(i64),
				}, cont, nil
			case reflect.Int64:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("int too large"))
				}
				return &sb.Token{
					Kind:  sb.KindInt64,
					Value: i64,
				}, cont, nil

			case reflect.Uint:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("uint too large"))
				}
				return &sb.Token{
					Kind:  sb.KindUint,
					Value: uint(u64),
				}, cont, nil
			case reflect.Uint8:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("uint too large"))
				}
				return &sb.Token{
					Kind:  sb.KindUint8,
					Value: uint8(u64),
				}, cont, nil
			case reflect.Uint16:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("uint too large"))
				}
				return &sb.Token{
					Kind:  sb.KindUint16,
					Value: uint16(u64),
				}, cont, nil
			case reflect.Uint32:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("uint too large"))
				}
				return &sb.Token{
					Kind:  sb.KindUint32,
					Value: uint32(u64),
				}, cont, nil
			case reflect.Uint64:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, nil, we(fmt.Errorf("uint too large"))
				}
				return &sb.Token{
					Kind:  sb.KindUint64,
					Value: u64,
				}, cont, nil

			default: // NOCOVER
				return nil, nil, we(fmt.Errorf("not an integer: %v", *t))
			}

		case starlark.Float:
			if t == nil {
				return &sb.Token{
					Kind:  sb.KindFloat64,
					Value: float64(value),
				}, cont, nil
			}
			switch (*t).Kind() {
			case reflect.Float32:
				return &sb.Token{
					Kind:  sb.KindFloat32,
					Value: float32(value),
				}, cont, nil
			case reflect.Float64:
				return &sb.Token{
					Kind:  sb.KindFloat64,
					Value: float64(value),
				}, cont, nil
			default: // NOCOVER
				return nil, nil, we(fmt.Errorf("not a float: %v", *t))
			}

		case starlark.String:
			return &sb.Token{
				Kind:  sb.KindString,
				Value: string(value),
			}, cont, nil

		case *starlark.List:
			iter := value.Iterate()
			var elemType *reflect.Type
			if t != nil {
				if kind := (*t).Kind(); kind != reflect.Slice && kind != reflect.Array { // NOCOVER
					return nil, nil, fmt.Errorf("not a list: %v", *t)
				}
				e := (*t).Elem()
				elemType = &e
			}
			return &sb.Token{
					Kind: sb.KindArray,
				}, marshalIter(
					iter,
					elemType,
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: sb.KindArrayEnd,
						}, cont, nil
					},
				), nil

		case starlark.Tuple:
			var types []reflect.Type
			if t != nil {
				if (*t).Kind() != reflect.Func { // NOCOVER
					return nil, nil, fmt.Errorf("not a tuple: %v", *t)
				}
				//TODO cache result
				for i, n := 0, (*t).NumOut(); i < n; i++ {
					types = append(types, (*t).Out(i))
				}
			}
			return &sb.Token{
					Kind: sb.KindTuple,
				}, marshalValues(
					value,
					types,
					func() (*sb.Token, Src, error) {
						return &sb.Token{
							Kind: sb.KindTupleEnd,
						}, cont, nil
					},
				), nil

		case *starlark.Dict:
			iter := value.Iterate()
			startKind := sb.KindMap
			endKind := sb.KindMapEnd
			if t != nil && (*t).Kind() != reflect.Map {
				startKind = sb.KindObject
				endKind = sb.KindObjectEnd
			}
			return &sb.Token{
					Kind: startKind,
				}, marshalDict(
					value,
					iter,
					t,
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: endKind,
						}, cont, nil
					},
				), nil

		case *starlark.Set:
			iter := value.Iterate()
			return &sb.Token{
					Kind: sb.KindMap,
				}, marshalSet(
					value,
					iter,
					t,
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: sb.KindMapEnd,
						}, cont, nil
					},
				), nil

		default: // NOCOVER
			panic(fmt.Errorf("unsupported type %T", value))

		}

	}
}

func marshalIter(
	iter starlark.Iterator,
	t *reflect.Type,
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		var value starlark.Value
		if iter.Next(&value) {
			return MarshalValue(value, t, marshalIter(iter, t, cont))()
		}
		return nil, cont, nil
	}
}

func marshalValues(
	values []starlark.Value,
	types []reflect.Type,
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		if len(values) == 0 {
			return nil, cont, nil
		}
		if len(types) > 0 {
			return MarshalValue(values[0], &types[0], marshalValues(values[1:], types[1:], cont))()
		}
		return MarshalValue(values[0], nil, marshalValues(values[1:], nil, cont))()
	}
}

func marshalDict(
	dict *starlark.Dict,
	iter starlark.Iterator,
	t *reflect.Type,
	cont Src,
) Src {
	var keyType *reflect.Type
	var valueType *reflect.Type
	if t != nil && (*t).Kind() == reflect.Map {
		kt := (*t).Key()
		keyType = &kt
		vt := (*t).Elem()
		valueType = &vt
	}
	return func() (*sb.Token, Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				keyType,
				func() (*sb.Token, Src, error) {
					value, ok, err := dict.Get(key)
					if err != nil { // NOCOVER
						return nil, nil, err
					}
					if !ok { // NOCOVER
						panic("impossible")
					}
					//TODO get type from struct field
					return MarshalValue(
						value,
						valueType,
						marshalDict(dict, iter, t, cont),
					)()
				},
			)()
		}
		return nil, cont, nil
	}
}

func marshalSet(
	set *starlark.Set,
	iter starlark.Iterator,
	t *reflect.Type,
	cont Src,
) Src {
	var keyType *reflect.Type
	if t != nil && (*t).Kind() == reflect.Map {
		kt := (*t).Key()
		keyType = &kt
	}
	return func() (*sb.Token, Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				keyType,
				func() (*sb.Token, Src, error) {
					return &sb.Token{
						Kind:  sb.KindBool,
						Value: true,
					}, marshalSet(set, iter, t, cont), nil
				},
			)()
		}
		return nil, cont, nil
	}
}
