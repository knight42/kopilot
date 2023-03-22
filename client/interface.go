package client

import "context"

type Client interface {
	CompletionRequest(ctx context.Context, prompt string) (string, error)
}
