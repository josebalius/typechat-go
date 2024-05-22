package typechat

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

type mockModelClient struct {
	response string
	err      error
}

func (m mockModelClient) Do(ctx context.Context, prompt []Message) (response string, err error) {
	return m.response, m.err
}

func TestTypeChat(t *testing.T) {
	type Result struct {
		Sentiment string `json:"sentiment"`
	}

	t.Run("it should generate the prompt and return the result", func(t *testing.T) {
		ctx := context.Background()
		m := mockModelClient{
			response: `{"sentiment": "positive"}`,
		}
		p := NewPrompt[Result](m, "That game was awesome!")
		result, err := p.Execute(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.Sentiment != "positive" {
			t.Errorf("Expected positive, got %v", result.Sentiment)
		}
	})

	t.Run("it should configure retries", func(t *testing.T) {
		p := NewPrompt[Result](nil, "", PromptRetries[Result](5))
		if p.retries != 5 {
			t.Errorf("Expected 5 retries, got %v", p.retries)
		}
	})

	t.Run("it should set a default number of retries", func(t *testing.T) {
		p := NewPrompt[Result](nil, "")
		if p.retries != 1 {
			t.Errorf("Expected 3 retries, got %v", p.retries)
		}
	})

	t.Run("it should accept an API interface and return a program", func(t *testing.T) {
		ctx := context.Background()
		type API interface {
			Step1(name string) (string, error)
			Step2(value int) error
		}
		program := Program{
			Steps: []FunctionCall{
				{
					Name: "Step1",
					Args: []any{"name"},
				},
				{
					Name: "Step2",
					Args: []any{2},
				},
			},
		}
		b, err := json.Marshal(program)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		m := mockModelClient{
			response: string(b),
		}
		p := NewPrompt[API](m, "")
		result, err := p.CreateProgram(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result.Steps) != 2 {
			t.Errorf("Expected 2 steps, got %v", len(result.Steps))
		}
		if result.Steps[0].Name != "Step1" {
			t.Errorf("Expected Step1, got %v", result.Steps[0].Name)
		}
		if result.Steps[1].Name != "Step2" {
			t.Errorf("Expected Step2, got %v", result.Steps[1].Name)
		}
	})
}

func TestProgram(t *testing.T) {
	t.Run("serialization of Program instance to JSON", func(t *testing.T) {
		program := Program{
			Steps: []FunctionCall{
				{Name: "Function1", Args: []interface{}{"arg1", 123}},
				{Name: "Function2", Args: []interface{}{"arg2", true}},
			},
		}
		expectedJSON := `{"Steps":[{"Name":"Function1","Args":["arg1",123]},{"Name":"Function2","Args":["arg2",true]}]}`
		bytes, err := json.Marshal(program)
		if err != nil {
			t.Fatalf("Failed to serialize Program: %v", err)
		}
		if string(bytes) != expectedJSON {
			t.Errorf("Serialized Program did not match expected JSON. Got: %s, Want: %s", string(bytes), expectedJSON)
		}
	})

	t.Run("deserialization of JSON to Program instance", func(t *testing.T) {
		jsonInput := `{"Steps":[{"Name":"Function1","Args":["arg1",123]},{"Name":"Function2","Args":["arg2",true]}]}`
		expectedProgram := Program{
			Steps: []FunctionCall{
				{Name: "Function1", Args: []interface{}{"arg1", 123}},
				{Name: "Function2", Args: []interface{}{"arg2", true}},
			},
		}
		var program Program
		err := json.Unmarshal([]byte(jsonInput), &program)
		if err != nil {
			t.Fatalf("Failed to deserialize JSON to Program: %v", err)
		}
		if !reflect.DeepEqual(program, expectedProgram) {
			t.Errorf("Deserialized Program did not match expected. Got: %+v, Want: %+v", program, expectedProgram)
		}
	})

	t.Run("Program with various FunctionCall configurations", func(t *testing.T) {
		tests := []struct {
			name     string
			program  Program
			expected string
		}{
			{
				name: "empty Program",
				program: Program{
					Steps: []FunctionCall{},
				},
				expected: `{"Steps":[]}`,
			},
			{
				name: "single FunctionCall with no args",
				program: Program{
					Steps: []FunctionCall{{Name: "Function1"}},
				},
				expected: `{"Steps":[{"Name":"Function1","Args":null}]}`,
			},
			{
				name: "multiple FunctionCalls with mixed args",
				program: Program{
					Steps: []FunctionCall{
						{Name: "Function1", Args: []interface{}{"arg1", 123}},
						{Name: "Function2", Args: []interface{}{}},
						{Name: "Function3", Args: []interface{}{"arg3", false}},
					},
				},
				expected: `{"Steps":[{"Name":"Function1","Args":["arg1",123]},{"Name":"Function2","Args":[]},{"Name":"Function3","Args":["arg3",false]}]}`,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bytes, err := json.Marshal(test.program)
				if err != nil {
					t.Fatalf("Failed to serialize Program: %v", err)
				}
				if string(bytes) != test.expected {
					t.Errorf("Serialized Program did not match expected JSON. Got: %s, Want: %s", string(bytes), test.expected)
				}
			})
		}
	})
}
