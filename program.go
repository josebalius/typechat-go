package typechat

import (
	"fmt"
	"reflect"
	"strings"
)

type Program struct {
	Steps []FunctionCall
}

type FunctionCall struct {
	Name string
	Args []interface{}
}

const (
	programSchemaInstructions = `You are a service that translates user requests into programs represented as JSON 
using the following Go definitions:`

	programPromptInstructions = `The following is the user request translated into a JSON object with 2 spaces of 
indentation and no properties with the value undefined:`
)

type program[T any] struct {
	input    string
	messages []Message
}

func newProgram[T any](i string) *program[T] {
	return &program[T]{input: i}
}

func (b *program[T]) prompt() ([]Message, error) {
	if b.messages != nil {
		return b.messages, nil
	}

	schema := new(T)
	schemaElem := reflect.TypeOf(schema).Elem()
	def, err := interfaceDef(schemaElem)
	if err != nil {
		return nil, fmt.Errorf("failed to get definition of schema: %w", err)
	}

	schemaPrompt, err := b.schema(def)
	if err != nil {
		return nil, fmt.Errorf("failed to build schema: %w", err)
	}

	b.messages = append(b.messages, newSystemMessage(schemaPrompt))
	b.messages = append(b.messages, newUserMessage(b.userMessage()))
	b.messages = append(b.messages, newSystemMessage(b.instructions()))

	return b.messages, nil
}

func (b *program[T]) userMessage() string {
	var sb strings.Builder
	sb.WriteString(newline("The following is a user request:"))
	sb.WriteString(newline(b.input))

	return sb.String()
}

func (b *program[T]) instructions() string {
	return programPromptInstructions
}

func (b *program[T]) schema(def string) (string, error) {
	var sb strings.Builder
	sb.WriteString(newline("A program consists of a sequence of function calls that are evaluated in order."))
	sb.WriteString(newline(programSchemaInstructions))

	_, programDef, err := structDef(reflect.TypeOf(Program{}))
	if err != nil {
		return "", err
	}
	sb.WriteString(programDef)

	sb.WriteString(newline("The programs can call functions from the API defined in the following Go definitions:"))
	sb.WriteString(def)

	return sb.String(), nil
}
