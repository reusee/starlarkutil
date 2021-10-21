package starlarkutil

import (
	"fmt"

	"github.com/reusee/sb"
	"go.starlark.net/starlark"
)

func Marshal(value starlark.Value, cont Src) *Src {
	src := MarshalValue(value, cont)
	return &src
}

func MarshalValue(value starlark.Value, cont Src) Src {
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
			i64, ok := value.Int64()
			if !ok {
				return nil, nil, we(fmt.Errorf("int too large"))
			}
			return &sb.Token{
				Kind:  sb.KindInt,
				Value: int(i64),
			}, cont, nil

		case starlark.Float:
			return &sb.Token{
				Kind:  sb.KindFloat64,
				Value: float64(value),
			}, cont, nil

		case starlark.String:
			return &sb.Token{
				Kind:  sb.KindString,
				Value: string(value),
			}, cont, nil

		case *starlark.List:
			iter := value.Iterate()
			return &sb.Token{
					Kind: sb.KindArray,
				}, marshalIter(
					iter,
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: sb.KindArrayEnd,
						}, cont, nil
					},
				), nil

		case starlark.Tuple:
			return &sb.Token{
					Kind: sb.KindTuple,
				}, marshalValues(
					value,
					func() (*sb.Token, Src, error) {
						return &sb.Token{
							Kind: sb.KindTupleEnd,
						}, cont, nil
					},
				), nil

		case *starlark.Dict:
			iter := value.Iterate()
			return &sb.Token{
					Kind: sb.KindMap,
				}, marshalDict(
					value,
					iter,
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: sb.KindMapEnd,
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
					func() (*sb.Token, Src, error) {
						iter.Done()
						return &sb.Token{
							Kind: sb.KindMapEnd,
						}, cont, nil
					},
				), nil

		default:
			panic(fmt.Errorf("unsupported type %T", value))

		}

		return nil, nil, nil
	}
}

func marshalIter(
	iter starlark.Iterator,
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		var value starlark.Value
		if iter.Next(&value) {
			return MarshalValue(value, marshalIter(iter, cont))()
		}
		return nil, cont, nil
	}
}

func marshalValues(
	values []starlark.Value,
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		if len(values) == 0 {
			return nil, cont, nil
		}
		return MarshalValue(values[0], marshalValues(values[1:], cont))()
	}
}

func marshalDict(
	dict *starlark.Dict,
	iter starlark.Iterator,
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				func() (*sb.Token, Src, error) {
					value, ok, err := dict.Get(key)
					if err != nil {
						return nil, nil, err
					}
					if !ok {
						panic("impossible")
					}
					return MarshalValue(
						value,
						marshalDict(dict, iter, cont),
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
	cont Src,
) Src {
	return func() (*sb.Token, Src, error) {
		var key starlark.Value
		if iter.Next(&key) {
			return MarshalValue(
				key,
				func() (*sb.Token, Src, error) {
					return &sb.Token{
						Kind:  sb.KindBool,
						Value: true,
					}, marshalSet(set, iter, cont), nil
				},
			)()
		}
		return nil, cont, nil
	}
}
