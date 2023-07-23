package typechat

import (
	"context"
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
}
