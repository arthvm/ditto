package gemini

import (
	"context"
	"sync"

	"github.com/arthvm/ditto/internal/llm"
	"google.golang.org/genai"
)

type Model = string

const (
	GeminiFlash     Model = "gemini-2.5-flash"
	GeminiFlashLite Model = "gemini-2.5-flash-lite"
	GeminiPro       Model = "gemini-2.5-pro"
)

type provider struct {
	model      Model
	clientOnce sync.Once
	client     *genai.Client
	clientErr  error
}

func init() {
	llm.Register("gemini", &provider{
		model: GeminiPro,
	})

	llm.Register("gemini-flash", &provider{
		model: GeminiFlash,
	})

	llm.Register("gemini-flash-lite", &provider{
		model: GeminiFlashLite,
	})
}

func (p *provider) getClient(ctx context.Context) (*genai.Client, error) {
	p.clientOnce.Do(func() {
		p.client, p.clientErr = genai.NewClient(ctx, nil)
	})

	return p.client, p.clientErr
}
