package typechat

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	userRequestSchemaInstructions = `You are a service that translates user requests into JSON objects of type %s 
according to the following Go definitions:`

	userRequestPromptInstructions = `The following is the user request translated into a JSON object with 2 spaces of 
indentation and no properties with the value undefined:`
)

type userRequest[T any] struct {
	builder *builder[T]
	built   string
}

func newUserRequest[T any](b *builder[T]) *userRequest[T] {
	return &userRequest[T]{builder: b}
}

func (b *userRequest[T]) string() (string, error) {
	if b.built != "" {
		return b.built, nil
	}

	var sb strings.Builder
	var schema T
	name, def, err := nameDef(reflect.TypeOf(schema))
	if err != nil {
		return "", err
	}

	writeLine(sb, b.schema(name, def))
	writeLine(sb, b.prompt())
	b.built = sb.String()

	return b.built, nil
}

func (b *userRequest[T]) prompt() string {
	var sb strings.Builder
	writeLine(sb, "The following is a user request:")
	writeLine(sb, b.builder.input)
	writeLine(sb, userRequestPromptInstructions)

	return sb.String()
}

func (b *userRequest[T]) schema(name, def string) string {
	var sb strings.Builder
	writeLine(sb, fmt.Sprintf(userRequestSchemaInstructions, name))
	writeLine(sb, def)

	return sb.String()
}
