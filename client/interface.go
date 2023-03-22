package client

import "context"

type Client interface {
	CreateCompletion(ctx context.Context, prompt string) (string, error)
}
