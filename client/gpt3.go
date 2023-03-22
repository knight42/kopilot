package client

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type GPT3Option struct {
	MaxToken int
}

type gpt3Client struct {
	client *openai.Client
	model  string
}

func NewGPT3Client(token string, options *GPT3Option) Client {
	c := openai.NewClient(token)
	return &gpt3Client{
		client: c,
		model:  openai.GPT3Ada,
	}
}

func (c *gpt3Client) CompletionRequest(ctx context.Context, prompt string) (string, error) {
	req := openai.CompletionRequest{
		Model:  c.model,
		Prompt: prompt,
	}
	resp, err := c.client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
