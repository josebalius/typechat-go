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
	input string
	built   string
}

func newUserRequest[T any](i string) *userRequest[T] {
  return &userRequest[T]{input: i}
}

func (b *userRequest[T]) string() (string, error) {
	if b.built != "" {
		return b.built, nil
	}

	var sb strings.Builder
	var schema T
	name, def, err := structDef(reflect.TypeOf(schema))
	if err != nil {
		return "", err
	}

	sb.WriteString(newline(b.schema(name, def)))
	sb.WriteString(newline(b.prompt()))
	b.built = sb.String()

	return b.built, nil
}

func (b *userRequest[T]) prompt() string {
	var sb strings.Builder
	sb.WriteString(newline("The following is a user request:"))
	sb.WriteString(newline(b.input))
	sb.WriteString(newline(userRequestPromptInstructions))

	return sb.String()
}

func (b *userRequest[T]) schema(name, def string) string {
	var sb strings.Builder
	sb.WriteString(newline(fmt.Sprintf(userRequestSchemaInstructions, name)))
	sb.WriteString(newline(def))

	return sb.String()
}
