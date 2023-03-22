package client

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPTOption struct {
	MaxToken int
	Lang     string
}

type chatGPTClient struct {
	client *openai.Client
	model  string
	lang   string
}

func NewChatGPTClient(token string, options ChatGPTOption) Client {
	c := openai.NewClient(token)
	return &chatGPTClient{
		client: c,
		model:  openai.GPT3Dot5Turbo,
		lang:   options.Lang,
	}
}

func (c *chatGPTClient) CreateCompletion(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
