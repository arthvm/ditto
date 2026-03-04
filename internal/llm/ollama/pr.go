package ollama

import (
	"context"

	"github.com/arthvm/ditto/internal/llm"
)

// TODO: Add support for PR creation on ollama models
func (p *provider) GeneratePR(
	ctx context.Context,
	params llm.GeneratePRParams,
) (string, error) {
	return "", llm.ErrNoSupport
}
