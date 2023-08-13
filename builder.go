package typechat

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type promptBuilder interface {
	prompt() ([]Message, error)
}

type builder[T any] struct {
	input string
	pt    promptType
	pb    promptBuilder
}

func newBuilder[T any](t promptType, input string) (*builder[T], error) {
	b := &builder[T]{
		input: input,
		pt:    t,
	}

	var pb promptBuilder
	switch t {
	case promptUserRequest:
		pb = newUserRequest[T](input)
	case promptProgram:
		pb = newProgram[T](input)
	default:
		return nil, fmt.Errorf("unknown prompt type %s", t)
	}
	b.pb = pb

	return b, nil
}

func (b *builder[T]) prompt() ([]Message, error) {
	return b.pb.prompt()
}

func (b *builder[T]) repair(resp string, err error) ([]Message, error) {
	msgs, err := b.pb.prompt()
	if err != nil {
		return nil, err
	}

	msgs = append(msgs, newAssistantMessage(resp))

	var sb strings.Builder
	if b.pt == promptUserRequest {
		sb.WriteString(newline("The JSON object is invalid for the following reason:"))
		sb.WriteString(newline(err.Error()))
		sb.WriteString(newline("The following is a revised JSON object:"))
	} else {
		sb.WriteString(newline("The JSON program object is invalid for the following reason:"))
		sb.WriteString(newline(err.Error()))
		sb.WriteString(newline("The following is a revised JSON program object:"))
	}

	msgs = append(msgs, newSystemMessage(sb.String()))

	return msgs, nil
}

func newline(s string) string {
	return fmt.Sprintf("%s\n", s)
}

func interfaceDef(t reflect.Type) (string, error) {
	if t.Kind() != reflect.Interface {
		return "", errors.New("top-level type must be an interface")
	}

	var methods strings.Builder
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		name := method.Name
		methodParts := []string{name}

		var args []string
		for j := 0; j < method.Type.NumIn(); j++ {
			in := method.Type.In(j)
			args = append(args, in.Name())
		}
		methodParts = append(methodParts, fmt.Sprintf("(%s)", strings.Join(args, ", ")))

		var returns []string
		for j := 0; j < method.Type.NumOut(); j++ {
			out := method.Type.Out(j)
			returns = append(returns, out.Name())
		}

		if len(returns) > 0 {
			methodParts = append(methodParts, fmt.Sprintf(" (%s)", strings.Join(returns, ", ")))
		}

		methods.WriteString(fmt.Sprintf("\t%s\n", strings.Join(methodParts, "")))
	}
	decl := fmt.Sprintf("type %s interface {\n%s}\n", t.Name(), methods.String())

	return decl, nil
}

func structDef(t reflect.Type) (string, string, error) {
	if t.Kind() != reflect.Struct {
		return "", "", errors.New("top-level type must be a struct")
	}

	name, decl, err := typeDecls(t)
	if err != nil {
		return "", "", err
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

		if disallowedField(field.Type, false) {
			return "", "", fmt.Errorf("field %s has disallowed type %s", field.Name, kind)
		}

		var structTags []string
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			structTags = append(structTags, fmt.Sprintf("json:\"%s\"", jsonTag))
		}

		if descriptionTag := field.Tag.Get("description"); descriptionTag != "" {
			structTags = append(structTags, fmt.Sprintf("description:\"%s\"", descriptionTag))
		}

		var structTag string
		if len(structTags) > 0 {
			structTag = fmt.Sprintf(" `%s`", strings.Join(structTags, " "))
		}

		name := field.Name
		typName := kind.String()
		if compositeField(kind) {
			n, decl, err := typeDecls(field.Type)
			if err != nil {
				return "", "", fmt.Errorf("field %s: %w", name, err)
			}
			typName = n
			decls.WriteString(decl)
		}

		fields.WriteString(fmt.Sprintf("\t%s %s%s\n", name, typName, structTag))
	}
	decls.WriteString(fmt.Sprintf("type %s struct {\n%s}\n", name, fields.String()))

	return name, decls.String(), nil
}

func typeSliceArrayDecl(t reflect.Type) (string, string, error) {
	var decls strings.Builder

	elem := t.Elem()
	name := elem.Name()
	kind := elem.Kind()

	if disallowedField(elem, true) {
		return "", "", fmt.Errorf("slice element has disallowed type %s", kind)
	}

	if compositeField(kind) {
		n, decl, err := typeDecls(elem)
		if err != nil {
			return "", "", err
		}
		name = n
		decls.WriteString(decl)
	}

	if kind == reflect.Interface {
		// empty interface case aka []any or []interface{}
		name = "interface{}"
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

	if disallowedField(key, false) {
		return "", "", fmt.Errorf("map key %s has disallowed type %s", keyName, keyKind)
	}

	if compositeField(keyKind) {
		return "", "", fmt.Errorf("map key %s has composite type %s", keyName, keyKind)
	}

	if disallowedField(value, false) {
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

func disallowedField(t reflect.Type, allowEmptyInterface bool) bool {
	k := t.Kind()
	if k == reflect.Complex64 ||
		k == reflect.Complex128 ||
		k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Pointer ||
		k == reflect.UnsafePointer {
		return true
	}

	if k == reflect.Interface {
		if t.NumMethod() == 0 && allowEmptyInterface {
			return false
		}

		return true
	}

	return false
}
