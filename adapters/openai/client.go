package openai

import (
	"context"
	"errors"

	"github.com/josebalius/typechat-go"
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

func openaiRole(r typechat.Role) (string, error) {
	var role string
	switch r {
	case typechat.RoleUser:
		role = openai.ChatMessageRoleUser
	case typechat.RoleSystem:
		role = openai.ChatMessageRoleSystem
	case typechat.RoleAssistant:
		role = openai.ChatMessageRoleAssistant
	default:
		return "", errors.New("invalid role")
	}

	return role, nil
}

func (c *Client) Do(ctx context.Context, prompt []typechat.Message) (string, error) {
	var messages []openai.ChatCompletionMessage
	for _, m := range prompt {
		role, err := openaiRole(m.Role)
		if err != nil {
			return "", err
		}

		msg := openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		}
		messages = append(messages, msg)
	}

	params := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
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
