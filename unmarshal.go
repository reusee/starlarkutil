package starlarkutil

import (
	"fmt"
	"io"

	"github.com/reusee/sb"
	"go.starlark.net/starlark"
)

func Unmarshal(target *starlark.Value, cont Sink) Sink {
	return func(token *sb.Token) (Sink, error) {
		if token.Invalid() { // NOCOVER
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(sb.UnmarshalError)
		}

		if target == nil { // NOCOVER
			var v starlark.Value
			target = &v
		}

		switch token.Kind {

		case sb.KindNil:

		case sb.KindBool:
			*target = starlark.Bool(token.Value.(bool))

		case sb.KindInt:
			*target = starlark.MakeInt(token.Value.(int))
		case sb.KindInt8:
			*target = starlark.MakeInt64(int64(token.Value.(int8)))
		case sb.KindInt16:
			*target = starlark.MakeInt64(int64(token.Value.(int16)))
		case sb.KindInt32:
			*target = starlark.MakeInt64(int64(token.Value.(int32)))
		case sb.KindInt64:
			*target = starlark.MakeInt64(token.Value.(int64))

		case sb.KindUint:
			*target = starlark.MakeUint(token.Value.(uint))
		case sb.KindUint8:
			*target = starlark.MakeUint64(uint64(token.Value.(uint8)))
		case sb.KindUint16:
			*target = starlark.MakeUint64(uint64(token.Value.(uint16)))
		case sb.KindUint32:
			*target = starlark.MakeUint64(uint64(token.Value.(uint32)))
		case sb.KindUint64:
			*target = starlark.MakeUint64(token.Value.(uint64))

		case sb.KindFloat32:
			*target = starlark.Float(token.Value.(float32))
		case sb.KindFloat64:
			*target = starlark.Float(token.Value.(float64))

		case sb.KindString:
			*target = starlark.String(token.Value.(string))

		case sb.KindBytes:
			*target = starlark.Bytes(token.Value.([]byte))

		case sb.KindArray:
			list := starlark.NewList(nil)
			return UnmarshalArray(
				list,
				func(token *sb.Token) (Sink, error) {
					*target = list
					return cont.Sink(token)
				},
			), nil

		case sb.KindObject, sb.KindMap:
			dict := starlark.NewDict(0)
			endKind := sb.KindObjectEnd
			if token.Kind == sb.KindMap {
				endKind = sb.KindMapEnd
			}
			return UnmarshalDict(
				endKind,
				dict,
				func(token *sb.Token) (Sink, error) {
					*target = dict
					return cont.Sink(token)
				},
			), nil

		case sb.KindTuple:
			var tuple starlark.Tuple
			return UnmarshalTuple(
				&tuple,
				func(token *sb.Token) (Sink, error) {
					*target = tuple
					return cont.Sink(token)
				},
			), nil

		default: // NOCOVER
			panic(fmt.Errorf("unknown token: %+v", token))
		}

		return cont, nil
	}
}

func UnmarshalArray(list *starlark.List, cont Sink) Sink {
	return func(token *sb.Token) (Sink, error) {
		if token.Invalid() { // NOCOVER
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(sb.UnmarshalError)
		}
		if token.Kind == sb.KindArrayEnd {
			return cont, nil
		}
		var elem starlark.Value
		return Unmarshal(
			&elem,
			func(token *sb.Token) (Sink, error) {
				list.Append(elem)
				return UnmarshalArray(list, cont).Sink(token)
			},
		).Sink(token)
	}
}

func UnmarshalDict(
	endKind sb.Kind,
	dict *starlark.Dict,
	cont Sink,
) Sink {
	return func(token *sb.Token) (Sink, error) {
		if token.Invalid() { // NOCOVER
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(sb.UnmarshalError)
		}
		if token.Kind == endKind {
			return cont, nil
		}
		var key starlark.Value
		return Unmarshal(
			&key,
			func(token *sb.Token) (Sink, error) {
				var value starlark.Value
				return Unmarshal(
					&value,
					func(token *sb.Token) (Sink, error) {
						dict.SetKey(key, value)
						return UnmarshalDict(endKind, dict, cont).Sink(token)
					},
				).Sink(token)
			},
		).Sink(token)
	}
}

func UnmarshalTuple(
	tuple *starlark.Tuple,
	cont Sink,
) Sink {
	return func(token *sb.Token) (Sink, error) {
		if token.Invalid() { // NOCOVER
			return nil, we.With(
				io.ErrUnexpectedEOF,
			)(sb.UnmarshalError)
		}
		if token.Kind == sb.KindTupleEnd {
			return cont, nil
		}
		var elem starlark.Value
		return Unmarshal(
			&elem,
			func(token *sb.Token) (Sink, error) {
				*tuple = append(*tuple, elem)
				return UnmarshalTuple(tuple, cont).Sink(token)
			},
		).Sink(token)
	}
}
