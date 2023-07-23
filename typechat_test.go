package typechat

import (
	"context"
	"encoding/json"
	"testing"
)

type mockModelClient struct {
	response []byte
	err      error
}

func (m mockModelClient) Do(ctx context.Context, prompt string) ([]byte, error) {
	return m.response, m.err
}

func TestTypeChat(t *testing.T) {
	type Result struct {
		Sentiment string `json:"sentiment"`
	}

	t.Run("it should generate the prompt and return the result", func(t *testing.T) {
		ctx := context.Background()
		m := mockModelClient{
			response: []byte(`{"sentiment": "positive"}`),
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
			response: b,
		}
		p := NewPrompt[API](m, "")
		result, err := p.ExecuteProgram(ctx)
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
