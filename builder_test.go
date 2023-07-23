package typechat

import (
	"reflect"
	"strings"
	"testing"
)

func assertNameDefOuptut(t *testing.T, def, expected string) {
	expected = strings.Trim(expected, "\n")
	def = strings.Trim(def, "\n")
	if def != expected {
		t.Errorf("expected def to be: \n%s\ngot: \n%s", expected, def)
	}
}

func TestNameDef(t *testing.T) {
	t.Run("it supports simple struct of primitive types", func(t *testing.T) {
		type Foo struct {
			A int
			B string
			C bool
		}

		name, def, err := nameDef(reflect.TypeOf(Foo{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		if name != "Foo" {
			t.Errorf("expected name to be Foo, got %s", name)
		}

		expected := `
type Foo struct {
	A int
	B string
	C bool
}`

		assertNameDefOuptut(t, def, expected)
	})

	t.Run("it supports json tags", func(t *testing.T) {
		type Foo struct {
			Field string `json:"field"`
		}

		_, def, err := nameDef(reflect.TypeOf(Foo{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		expected := `
type Foo struct {
	Field string ` + "`" + `json:"field"` + "`" + `
}`
		assertNameDefOuptut(t, def, expected)
	})

	t.Run("it supports structs", func(t *testing.T) {
		type Foo struct {
			A int
			B string
			C bool
		}

		type Bar struct {
			Foo Foo
		}

		name, def, err := nameDef(reflect.TypeOf(Bar{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		if name != "Bar" {
			t.Errorf("expected name to be Bar, got %s", name)
		}

		expected := `
type Foo struct {
	A int
	B string
	C bool
}
type Bar struct {
	Foo Foo
}`

		assertNameDefOuptut(t, def, expected)
	})

	t.Run("it supports slices", func(t *testing.T) {
		type Bar struct {
			C string
		}

		type Foo struct {
			A []int
			B []Bar
		}

		_, def, err := nameDef(reflect.TypeOf(Foo{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		expected := `
type Bar struct {
	C string
}
type Foo struct {
	A []int
	B []Bar
}`
		assertNameDefOuptut(t, def, expected)
	})

	t.Run("it supports maps", func(t *testing.T) {
		type Bar struct {
			C string
		}

		type Foo struct {
			A []int
			B map[string]Bar
		}

		_, def, err := nameDef(reflect.TypeOf(Foo{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		expected := `
type Bar struct {
	C string
}
type Foo struct {
	A []int
	B map[string]Bar
}`
		assertNameDefOuptut(t, def, expected)
	})

	t.Run("it supports nested maps", func(t *testing.T) {
		type Bar struct {
			C string
		}

		type Foo struct {
			A []int
			B map[string]map[string]Bar
		}

		_, def, err := nameDef(reflect.TypeOf(Foo{}))
		if err != nil {
			t.Fatalf("expected err to be nil, got %s", err)
		}

		expected := `
type Bar struct {
	C string
}
type Foo struct {
	A []int
	B map[string]map[string]Bar
}`
		assertNameDefOuptut(t, def, expected)
	})

}
