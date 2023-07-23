package typechat

import (
	"reflect"
	"strings"
)

const (
	programSchemaInstructions = `You are a service that translates user requests into programs represented as JSON 
using the following Go definitions:`

	programText = `// A program consists of a sequence of function calls that are evaluated in order.
type Program struct {
	Steps []FunctionCall
}

type FunctionCall struct {
	Name string
	Args []any
}`

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

func (b *program[T]) prompt() string {
	var sb strings.Builder
	writeLine(sb, "The following is a user request:")
	writeLine(sb, b.builder.input)
	writeLine(sb, programPromptInstructions)

	return sb.String()
}

func (b *program[T]) schema(name, def string) string {
	var sb strings.Builder
	writeLine(sb, programSchemaInstructions)
	writeLine(sb, programText)
	writeLine(sb, "The programs can call functions from the API defined in the following Go definitions:")
	writeLine(sb, def)

	return sb.String()
}
