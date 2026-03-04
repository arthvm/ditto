package llm

import "context"

// Provider generates text from a system prompt and user prompt.
type Provider interface {
	Generate(ctx context.Context, system, user string) (string, error)
}
