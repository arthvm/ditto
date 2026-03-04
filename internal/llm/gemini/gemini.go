package gemini

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/genai"
)

// Provider implements workflow.Provider using the Google Gemini API.
// The underlying client is lazily initialized on first Generate call.
type Provider struct {
	model       string
	temperature float32
	clientOnce  sync.Once
	client      *genai.Client
	clientErr   error
}

// New creates a Gemini provider for the given model name
// (e.g. "gemini-2.5-flash", "gemini-2.5-pro").
// A temperature of 0 uses the model's default.
func New(model string, temperature float32) *Provider {
	return &Provider{model: model, temperature: temperature}
}

func (p *Provider) getClient(ctx context.Context) (*genai.Client, error) {
	p.clientOnce.Do(func() {
		p.client, p.clientErr = genai.NewClient(ctx, nil)
	})

	return p.client, p.clientErr
}

func (p *Provider) Generate(ctx context.Context, system, user string) (string, error) {
	client, err := p.getClient(ctx)
	if err != nil {
		return "", fmt.Errorf("generate client: %w", err)
	}

	cfg := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(system, genai.RoleUser),
	}
	if p.temperature != 0 {
		cfg.Temperature = &p.temperature
	}

	result, err := client.Models.GenerateContent(
		ctx,
		p.model,
		genai.Text(user),
		cfg,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
