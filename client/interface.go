package client

import (
	"context"
	"io"
)

type Client interface {
	CreateCompletion(ctx context.Context, prompt string, writer io.Writer) error
}
