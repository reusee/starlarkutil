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
	return func(token *sb.Token) (Src, error) {

		switch value := value.(type) {

		case starlark.NoneType:
			token.Kind = sb.KindNil
			return cont, nil

		case starlark.Bool:
			token.Kind = sb.KindBool
			token.Value = bool(value)
			return cont, nil

		case starlark.Bytes:
			token.Kind = sb.KindBytes
			token.Value = []byte(value)
			return cont, nil

		case starlark.Int:
			if t == nil {
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt
				token.Value = int(i64)
				return cont, nil
			}

			switch (*t).Kind() {
			case reflect.Int:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt
				token.Value = int(i64)
				return cont, nil
			case reflect.Int8:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt8
				token.Value = int8(i64)
				return cont, nil
			case reflect.Int16:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt16
				token.Value = int16(i64)
				return cont, nil
			case reflect.Int32:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt32
				token.Value = int32(i64)
				return cont, nil
			case reflect.Int64:
				i64, ok := value.Int64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("int too large"))
				}
				token.Kind = sb.KindInt64
				token.Value = i64
				return cont, nil

			case reflect.Uint:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("uint too large"))
				}
				token.Kind = sb.KindUint
				token.Value = uint(u64)
				return cont, nil
			case reflect.Uint8:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("uint too large"))
				}
				token.Kind = sb.KindUint8
				token.Value = uint8(u64)
				return cont, nil
			case reflect.Uint16:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("uint too large"))
				}
				token.Kind = sb.KindUint16
				token.Value = uint16(u64)
				return cont, nil
			case reflect.Uint32:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("uint too large"))
				}
				token.Kind = sb.KindUint32
				token.Value = uint32(u64)
				return cont, nil
			case reflect.Uint64:
				u64, ok := value.Uint64()
				if !ok { // NOCOVER
					return nil, we(fmt.Errorf("uint too large"))
				}
				token.Kind = sb.KindUint64
				token.Value = u64
				return cont, nil

			default: // NOCOVER
				return nil, we(fmt.Errorf("not an integer: %v", *t))
			}

		case starlark.Float:
			if t == nil {
				token.Kind = sb.KindFloat64
				token.Value = float64(value)
				return cont, nil
			}
			switch (*t).Kind() {
			case reflect.Float32:
				token.Kind = sb.KindFloat32
				token.Value = float32(value)
				return cont, nil
			case reflect.Float64:
				token.Kind = sb.KindFloat64
				token.Value = float64(value)
				return cont, nil
			default: // NOCOVER
				return nil, we(fmt.Errorf("not a float: %v", *t))
			}

		case starlark.String:
			token.Kind = sb.KindString
			token.Value = string(value)
			return cont, nil

		case *starlark.List:
			iter := value.Iterate()
			var elemType *reflect.Type
			if t != nil {
				if kind := (*t).Kind(); kind != reflect.Slice && kind != reflect.Array { // NOCOVER
					return nil, fmt.Errorf("not a list: %v", *t)
				}
				e := (*t).Elem()
				elemType = &e
			}
			token.Kind = sb.KindArray
			return marshalIter(
				iter,
				elemType,
				func(token *sb.Token) (Src, error) {
					iter.Done()
					token.Kind = sb.KindArrayEnd
					return cont, nil
				},
			), nil

		case starlark.Tuple:
			var types []reflect.Type
			if t != nil {
				if (*t).Kind() != reflect.Func { // NOCOVER
					return nil, fmt.Errorf("not a tuple: %v", *t)
				}
				//TODO cache result
				for i, n := 0, (*t).NumOut(); i < n; i++ {
					types = append(types, (*t).Out(i))
				}
			}
			token.Kind = sb.KindTuple
			return marshalValues(
				value,
				types,
				func(token *sb.Token) (Src, error) {
					token.Kind = sb.KindTupleEnd
					return cont, nil
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
			token.Kind = startKind
			return marshalDict(
				value,
				iter,
				t,
				func(token *sb.Token) (Src, error) {
					iter.Done()
					token.Kind = endKind
					return cont, nil
				},
			), nil

		case *starlark.Set:
			iter := value.Iterate()
			token.Kind = sb.KindMap
			return marshalSet(
				value,
				iter,
				t,
				func(token *sb.Token) (Src, error) {
					iter.Done()
					token.Kind = sb.KindMapEnd
					return cont, nil
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
	return func(token *sb.Token) (Src, error) {
		var value starlark.Value
		if iter.Next(&value) {
			return MarshalValue(value, t, marshalIter(iter, t, cont))(token)
		}
		return cont, nil
	}
}

func marshalValues(
	values []starlark.Value,
	types []reflect.Type,
	cont Src,
) Src {
	return func(token *sb.Token) (Src, error) {
		if len(values) == 0 {
			return cont, nil
		}
		if len(types) > 0 {
			return MarshalValue(values[0], &types[0], marshalValues(values[1:], types[1:], cont))(token)
		}
		return MarshalValue(values[0], nil, marshalValues(values[1:], nil, cont))(token)
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
	return func(token *sb.Token) (Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				keyType,
				func(token *sb.Token) (Src, error) {
					value, ok, err := dict.Get(key)
					if err != nil { // NOCOVER
						return nil, err
					}
					if !ok { // NOCOVER
						panic("impossible")
					}
					//TODO get type from struct field
					return MarshalValue(
						value,
						valueType,
						marshalDict(dict, iter, t, cont),
					)(token)
				},
			)(token)
		}
		return cont, nil
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
	return func(token *sb.Token) (Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				keyType,
				func(token *sb.Token) (Src, error) {
					token.Kind = sb.KindBool
					token.Value = true
					return marshalSet(set, iter, t, cont), nil
				},
			)(token)
		}
		return cont, nil
	}
}
