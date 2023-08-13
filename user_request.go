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
	input    string
	messages []Message
}

func newUserRequest[T any](i string) *userRequest[T] {
	return &userRequest[T]{input: i}
}

func (b *userRequest[T]) prompt() ([]Message, error) {
	if b.messages != nil {
		return b.messages, nil
	}

	var schema T
	name, def, err := structDef(reflect.TypeOf(schema))
	if err != nil {
		return nil, err
	}

	b.messages = append(b.messages, newSystemMessage(b.schema(name, def)))
	b.messages = append(b.messages, newUserMessage(b.userMessage()))
	b.messages = append(b.messages, newSystemMessage(b.instructions()))

	return b.messages, nil
}

func (b *userRequest[T]) userMessage() string {
	var sb strings.Builder
	sb.WriteString(newline("The following is a user request:"))
	sb.WriteString(newline(b.input))

	return sb.String()
}

func (b *userRequest[T]) instructions() string {
	return userRequestPromptInstructions
}

func (b *userRequest[T]) schema(name, def string) string {
	var sb strings.Builder
	sb.WriteString(newline(fmt.Sprintf(userRequestSchemaInstructions, name)))
	sb.WriteString(newline(def))

	return sb.String()
}
