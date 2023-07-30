package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	model  string
}

func NewClient(client *openai.Client, model string) *Client {
	return &Client{
		client: client,
		model:  model,
	}
}

func (c *Client) Do(ctx context.Context, prompt string) (string, error) {
	params := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := c.client.CreateChatCompletion(ctx, params)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
