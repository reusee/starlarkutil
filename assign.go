package starlarkutil

import (
	"reflect"

	"go.starlark.net/starlark"
)

func Assign(value starlark.Value, target reflect.Value) error {

	targetType := target.Type()
	targetKind := targetType.Kind()
	if targetKind != reflect.Ptr {
		return &InvalidTarget{
			Str:  value.String(),
			Type: targetType,
		}
	}
	valueType := targetType.Elem()
	valueKind := valueType.Kind()
	if valueKind != reflect.Interface {
		if target.IsNil() {
			target = reflect.New(valueType)
		}
	}

	switch valueKind {

	case reflect.Bool:
		target.Elem().Set(reflect.ValueOf(value.Truth() == starlark.True))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := starlark.NumberToInt(value)
		if err != nil {
			return we.With(
				WithInvalidValue(value.String()),
			)(err)
		}
		n, ok := i.Int64()
		if !ok {
			return &InvalidValue{
				Str: value.String(),
			}
		}
		target.Elem().SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := starlark.NumberToInt(value)
		if err != nil {
			return we.With(
				WithInvalidValue(value.String()),
			)(err)
		}
		n, ok := i.Uint64()
		if !ok {
			return &InvalidValue{
				Str: value.String(),
			}
		}
		target.Elem().SetUint(n)

	case reflect.Float32, reflect.Float64:
		f, ok := starlark.AsFloat(value)
		if !ok {
			return &InvalidValue{
				Str: value.String(),
			}
		}
		target.Elem().SetFloat(f)

	case reflect.Array:
		iter := starlark.Iterate(value)
		var v starlark.Value
		i := 0
		for iter.Next(&v) {
			if i >= target.Elem().Len() {
				return we.With(
					WithInvalidValue(value.String()),
				)(TooManyElements)
			}
			if err := Assign(v, target.Elem().Index(i).Addr()); err != nil {
				return err
			}
			i++
		}
		iter.Done()

	case reflect.Interface:
		switch value.Type() {

		case "int":
			var i int
			if err := Assign(value, reflect.ValueOf(&i)); err != nil {
				return err
			}
			target.Elem().Set(reflect.ValueOf(i))

		case "float":
			var f float64
			if err := Assign(value, reflect.ValueOf(&f)); err != nil { // NOCOVER
				return err
			}
			target.Elem().Set(reflect.ValueOf(f))

		case "string":
			var s string
			if err := Assign(value, reflect.ValueOf(&s)); err != nil { // NOCOVER
				return err
			}
			target.Elem().Set(reflect.ValueOf(s))

		case "list":
			var l []any
			if err := Assign(value, reflect.ValueOf(&l)); err != nil {
				return err
			}
			target.Elem().Set(reflect.ValueOf(l))

		case "dict":
			m := make(map[any]any)
			if err := Assign(value, reflect.ValueOf(&m)); err != nil {
				return err
			}
			target.Elem().Set(reflect.ValueOf(m))

		case "bool":
			var b bool
			if err := Assign(value, reflect.ValueOf(&b)); err != nil { // NOCOVER
				return err
			}
			target.Elem().Set(reflect.ValueOf(b))

		default:
			return &InvalidValue{
				Str: value.String(),
			}

		}

	case reflect.Map:
		iter := starlark.Iterate(value)
		var key starlark.Value
		m := value.(starlark.Mapping)
		for iter.Next(&key) {
			v, ok, err := m.Get(key)
			if err != nil { // NOCOVER
				return err
			}
			if !ok { // NOCOVER
				continue
			}
			var goValue any
			if err := Assign(v, reflect.ValueOf(&goValue)); err != nil {
				return err
			}
			var goKey any
			if err := Assign(key, reflect.ValueOf(&goKey)); err != nil {
				return err
			}
			if target.Elem().IsNil() {
				target.Elem().Set(reflect.MakeMap(valueType))
			}
			target.Elem().SetMapIndex(
				reflect.ValueOf(goKey),
				reflect.ValueOf(goValue),
			)
		}
		iter.Done()

	case reflect.Ptr:
		if target.Elem().IsNil() {
			t := reflect.New(valueType.Elem())
			if err := Assign(value, t); err != nil {
				return err
			}
			target.Elem().Set(t)
		} else {
			return Assign(value, target.Elem())
		}

	case reflect.Slice:
		iter := starlark.Iterate(value)
		var v starlark.Value
		slice := target.Elem()
		for iter.Next(&v) {
			elem := reflect.New(valueType.Elem())
			if err := Assign(v, elem); err != nil {
				return err
			}
			slice = reflect.Append(slice, elem.Elem())
		}
		iter.Done()
		target.Elem().Set(slice)

	case reflect.String:
		s, ok := starlark.AsString(value)
		if !ok {
			return &InvalidValue{
				Str: value.String(),
			}
		}
		target.Elem().SetString(s)

	case reflect.Struct:
		iter := starlark.Iterate(value)
		var key starlark.Value
		m := value.(starlark.Mapping)
		for iter.Next(&key) {
			name, ok := starlark.AsString(key)
			if !ok {
				continue
			}
			field := target.Elem().FieldByName(name)
			if !field.IsValid() {
				continue
			}
			v, ok, err := m.Get(key)
			if err != nil { // NOCOVER
				return err
			}
			if !ok { // NOCOVER
				continue
			}
			if err := Assign(v, field.Addr()); err != nil {
				return err
			}
		}
		iter.Done()

	default:
		return &InvalidTarget{
			Str:  value.String(),
			Type: reflect.TypeOf(target),
		}

	}

	return nil
}
