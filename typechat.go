package typechat

import (
	"context"
	"encoding/json"
	"fmt"
)

type client interface {
	Do(ctx context.Context, prompt string) (response []byte, err error)
}

// Prompt is a generic typechat prompt.
type Prompt[T any] struct {
	model  client
	prompt string

	retries int
}

type opt[T any] func(*Prompt[T])

// PromptRetries sets the number of times to retry parsing errors.
func PromptRetries[T any](retries int) opt[T] {
	return func(t *Prompt[T]) {
		t.retries = retries
	}
}

// NewPrompt creates a new Prompt[T] with the given modelClient, prompt and options.
func NewPrompt[T any](model client, prompt string, opts ...opt[T]) *Prompt[T] {
	t := &Prompt[T]{
		model:  model,
		prompt: prompt,
	}
	for _, opt := range opts {
		opt(t)
	}
	if t.retries <= 0 {
		t.retries = 1
	}

	return t
}

// Execute executes the user request prompt and parses the result. Parsing errors are retried up to Prompt.retries
// times.
func (p *Prompt[T]) Execute(ctx context.Context) (T, error) {
	var result T

	b, err := newBuilder[T](promptUserRequest, p.prompt)
	if err != nil {
		return result, fmt.Errorf("failed to create prompt builder: %w", err)
	}

	return p.exec(ctx, b)
}

// ExecuteProgram executes the program prompt and parses the result. Parsing errors are retried up to Prompt.retries
// times.
func (p *Prompt[T]) ExecuteProgram(ctx context.Context) (T, error) {
	var result T

	b, err := newBuilder[T](promptProgram, p.prompt)
	if err != nil {
		return result, fmt.Errorf("failed to create prompt builder: %w", err)
	}

	return p.exec(ctx, b)
}

func (p *Prompt[T]) exec(ctx context.Context, b *builder[T]) (T, error) {
	var result T
	prompt, err := b.string()
	if err != nil {
		return result, fmt.Errorf("failed to build prompt: %w", err)
	}

	for i := 0; i < p.retries; i++ {
		resp, err := p.model.Do(ctx, prompt)
		if err != nil {
			return result, err
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			prompt = b.repair(resp, err)
			continue
		}
	}

	return result, nil
}
