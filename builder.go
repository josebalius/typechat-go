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
	writeLine(sb, string(resp))
	writeLine(sb, "The JSON object is invalid for the following reason:")
	writeLine(sb, err.Error())
	writeLine(sb, "The following is a revised JSON object:")

	return sb.String()
}

func writeLine(b strings.Builder, s string) {
	b.WriteString(fmt.Sprintf("%s\n", s))
}

func nameDef(t reflect.Type) (string, string, error) {
	fields := ""

	var compositeFields []reflect.Type
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := field.Type.Kind()
		if disallowedField(kind) {
			return "", "", fmt.Errorf("field %s has disallowed type %s", field.Name, kind)
		}
		if compositeField(kind) {
			compositeFields = append(compositeFields, field.Type)
		}

		jsonTag := field.Tag.Get("json")
		var jsonTagDef string
		if jsonTag != "" {
			jsonTagDef = fmt.Sprintf(" `json:\"%s\"`", jsonTag)
		}

		fields += fmt.Sprintf("\t%s %s%s\n", field.Name, kind, jsonTagDef)
	}

	name := t.Name()
	return name, fmt.Sprintf("type %s struct {\n%s}", name, fields), nil
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
