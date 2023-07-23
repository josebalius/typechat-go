package typechat

import (
	"fmt"
	"reflect"
	"strings"
)

type promptBuilder interface {
	string() (string, error)
}

type builder[T any] struct {
	input string
	pb    promptBuilder
}

func newBuilder[T any](t promptType, input string) (*builder[T], error) {
	b := &builder[T]{
		input: input,
	}

	var pb promptBuilder
	switch t {
	case promptUserRequest:
		pb = newUserRequest[T](b)
	case promptProgram:
		pb = newProgram[T](b)
	default:
		return nil, fmt.Errorf("unknown prompt type %s", t)
	}
	b.pb = pb

	return b, nil
}

func (b *builder[T]) string() (string, error) {
	return b.pb.string()
}

func (b *builder[T]) repair(resp []byte, err error) string {
	var sb strings.Builder
	sb.WriteString(newline(string(resp)))
	sb.WriteString(newline("The JSON object is invalid for the following reason:"))
	sb.WriteString(newline(err.Error()))
	sb.WriteString(newline("The following is a revised JSON object:"))

	return sb.String()
}

func newline(s string) string {
	return fmt.Sprintf("%s\n", s)
}

func nameDef(t reflect.Type) (string, string, error) {
	name, decl, err := typeDecls(t)
	if err != nil {
		return "", "", fmt.Errorf("failed to define type: %w", err)
	}

	return name, decl, nil
}

func typeDecls(t reflect.Type) (string, string, error) {
	switch t.Kind() {
	case reflect.Struct:
		return typeStructDecl(t)
	case reflect.Slice, reflect.Array:
		return typeSliceArrayDecl(t)
	case reflect.Map:
		return typeMapDecl(t)
	default:
		return "", "", fmt.Errorf("unsupported type %s", t.Kind())
	}
}

func typeStructDecl(t reflect.Type) (string, string, error) {
	name := t.Name()

	var decls strings.Builder
	var fields strings.Builder

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := field.Type.Kind()

		if disallowedField(kind) {
			return "", "", fmt.Errorf("field %s has disallowed type %s", field.Name, kind)
		}

		jsonTag := field.Tag.Get("json")
		var jsonTagDef string
		if jsonTag != "" {
			jsonTagDef = fmt.Sprintf(" `json:\"%s\"`", jsonTag)
		}

		name := field.Name
		typName := kind.String()
		if compositeField(kind) {
			n, decl, err := typeDecls(field.Type)
			if err != nil {
				return "", "", fmt.Errorf("failed to define type: %w", err)
			}
			typName = n
			decls.WriteString(decl)
		}

		fields.WriteString(fmt.Sprintf("\t%s %s%s\n", name, typName, jsonTagDef))
	}

	decls.WriteString(fmt.Sprintf("type %s struct {\n%s}\n", name, fields.String()))
	return name, decls.String(), nil
}

func typeSliceArrayDecl(t reflect.Type) (string, string, error) {
	var decls strings.Builder

	elem := t.Elem()
	name := elem.Name()
	kind := elem.Kind()
	if disallowedField(kind) {
		return "", "", fmt.Errorf("slice element has disallowed type %s", kind)
	}
	if compositeField(kind) {
		n, decl, err := typeDecls(elem)
		if err != nil {
			return "", "", fmt.Errorf("failed to define type: %w", err)
		}
		name = n
		decls.WriteString(decl)
	}

	name = fmt.Sprintf("[]%s", name)
	return name, decls.String(), nil
}

func typeMapDecl(t reflect.Type) (string, string, error) {
	var decls strings.Builder

	key := t.Key()
	keyName := key.Name()
	keyKind := key.Kind()

	value := t.Elem()
	valueName := value.Name()
	valueKind := value.Kind()

	if disallowedField(keyKind) {
		return "", "", fmt.Errorf("map key %s has disallowed type %s", keyName, keyKind)
	}
	if compositeField(keyKind) {
		return "", "", fmt.Errorf("map key %s has composite type %s", keyName, keyKind)
	}
	if disallowedField(valueKind) {
		return "", "", fmt.Errorf("map value %s has disallowed type %s", valueName, valueKind)
	}
	if compositeField(valueKind) {
		n, decl, err := typeDecls(value)
		if err != nil {
			return "", "", fmt.Errorf("failed to define type: %w", err)
		}
		valueName = n
		decls.WriteString(decl)
	}

	name := fmt.Sprintf("map[%s]%s", keyName, valueName)
	return name, decls.String(), nil
}

func compositeField(k reflect.Kind) bool {
	return k == reflect.Map ||
		k == reflect.Slice ||
		k == reflect.Array ||
		k == reflect.Struct
}

func disallowedField(k reflect.Kind) bool {
	return k == reflect.Complex64 ||
		k == reflect.Complex128 ||
		k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Pointer ||
		k == reflect.UnsafePointer
}
