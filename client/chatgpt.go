package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
)

type chatGPTClient struct {
	client *openai.Client
	model  string
}

func NewChatGPTClient(token string, proxy string, proto string) Client {
	c := &openai.Client{}
	if proxy != "" && proto != "" {
		config := openai.DefaultConfig(token)
		config.HTTPClient.Transport = &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse(proto + "://" + proxy)
			},
		}
		c = openai.NewClientWithConfig(config)
	} else {
		c = openai.NewClient(token)
	}

	return &chatGPTClient{
		client: c,
		model:  openai.GPT3Dot5Turbo,
	}
}

func (c *chatGPTClient) CreateCompletion(ctx context.Context, prompt string, writer io.Writer, spinner Spinner) error {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	isSpinnerStopped := false
	isFirstMsg := true

	defer func() {
		if !isSpinnerStopped {
			spinner.Stop()
		}
	}()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		if isFirstMsg {
			spinner.Stop()
			isFirstMsg = false
			isSpinnerStopped = true
		}

		_, err = io.WriteString(writer, response.Choices[0].Delta.Content)
		if err != nil {
			return err
		}
	}
}
