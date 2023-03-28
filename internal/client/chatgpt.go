package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sashabaranov/go-openai"
)

type chatGPTClient struct {
	client *openai.Client
	model  string

	spinnerSuffix string
}

func NewChatGPTClient(token, spinnerSuffix string) Client {
	c := openai.NewClient(token)
	return &chatGPTClient{
		client: c,
		model:  openai.GPT3Dot5Turbo,

		spinnerSuffix: spinnerSuffix,
	}
}

func (c *chatGPTClient) CreateCompletion(ctx context.Context, prompt string, writer io.Writer) error {
	s := spinner.New(spinner.CharSets[16], 100*time.Millisecond)
	s.Suffix = c.spinnerSuffix
	s.Start()
	// it is okay to stop spinner twice
	defer s.Stop()

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

	isFirstMsg := true

	for {
		response, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(writer)
				return nil
			}
			return fmt.Errorf("receive stream: %w", err)
		}

		if isFirstMsg {
			s.Stop()
			isFirstMsg = false
		}

		_, err = io.WriteString(writer, response.Choices[0].Delta.Content)
		if err != nil {
			return err
		}
	}
}
