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

	if err := p.exec(ctx, b, &result); err != nil {
		return result, fmt.Errorf("failed to execute prompt: %w", err)
	}

	return result, nil
}

// ExecuteProgram executes the program prompt and parses the result. Parsing errors are retried up to Prompt.retries
// times.
func (p *Prompt[T]) ExecuteProgram(ctx context.Context) (Program, error) {
	var program Program

	b, err := newBuilder[T](promptProgram, p.prompt)
	if err != nil {
		return program, fmt.Errorf("failed to create prompt builder: %w", err)
	}

	if err := p.exec(ctx, b, &program); err != nil {
		return program, fmt.Errorf("failed to execute prompt: %w", err)
	}

	return program, nil
}

func (p *Prompt[T]) exec(ctx context.Context, b *builder[T], output any) error {
	prompt, err := b.string()
	fmt.Println(prompt)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}

	for i := 0; i < p.retries; i++ {
		resp, err := p.model.Do(ctx, prompt)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(resp, output); err != nil {
			prompt = b.repair(resp, err)
			continue
		}
	}

	return nil
}
