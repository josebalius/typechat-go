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
	builder *builder[T]
	built   string
}

func newProgram[T any](b *builder[T]) *program[T] {
	return &program[T]{builder: b}
}

func (b *program[T]) string() (string, error) {
	if b.built != "" {
		return b.built, nil
	}

	var sb strings.Builder
	schema := new(T)

	schemaElem := reflect.TypeOf(schema).Elem()
	def, err := interfaceDef(schemaElem)
	if err != nil {
		return "", fmt.Errorf("failed to get definition of schema: %w", err)
	}

	schemaPrompt, err := b.schema(def)
	if err != nil {
		return "", fmt.Errorf("failed to build schema: %w", err)
	}

	sb.WriteString(schemaPrompt)
	sb.WriteString(b.prompt())
	b.built = sb.String()

	return b.built, nil
}

func (b *program[T]) prompt() string {
	var sb strings.Builder
	sb.WriteString(newline("The following is a user request:"))
	sb.WriteString(newline(b.builder.input))
	sb.WriteString(newline(programPromptInstructions))

	return sb.String()
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
