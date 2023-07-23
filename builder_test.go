package typechat

import (
	"strings"
	"testing"
)

func TestTypeOf(t *testing.T) {
	t.Run("it supports simple struct of primitive types", func(t *testing.T) {
		type Foo struct {
			A int
			B string
			C bool
		}

		name, def := typeOf(Foo{})

		if name != "Foo" {
			t.Errorf("expected name to be Foo, got %s", name)
		}

		expected := `
type Foo struct {
	A int
	B string
	C bool
}`
		if def != strings.TrimPrefix(expected, "\n") {
			t.Errorf("expected def to be %s, got %s", expected, def)
		}
	})

	t.Run("it supports json tags", func(t *testing.T) {
		type Foo struct {
			Field string `json:"field"`
		}

		_, def := typeOf(Foo{})

		expected := `
type Foo struct {
	Field string ` + "`" + `json:"field"` + "`" + `
}`
		expected = strings.TrimPrefix(expected, "\n")
		if def != expected {
			t.Errorf("expected def to be: \n%s \ngot \n%s", expected, def)
		}
	})

	t.Run("it supports complex types", func(t *testing.T) {
		type Foo struct {
			A int
			B string
			C bool
		}

		type Bar struct {
			Foo Foo
		}

		name, def := typeOf(Bar{})

		if name != "Bar" {
			t.Errorf("expected name to be Bar, got %s", name)
		}

		expected := `
type Bar struct {
	Foo Foo
}`

		expected = strings.TrimPrefix(expected, "\n")
		if def != expected {
			t.Errorf("expected def to be %s, got %s", expected, def)
		}
	})

}
