package ollama

import (
	"context"

	"github.com/arthvm/ditto/internal/llm"
)

// TODO: Add suppport for PR creation on ollama models
func (p *provider) GeneratePr(
	ctx context.Context,
	params llm.GeneratePrParams,
) (string, error) {
	return "", llm.ErrNoSupport
}
