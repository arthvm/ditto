package gemini

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/genai"
)

type Provider struct {
	model      string
	clientOnce sync.Once
	client     *genai.Client
	clientErr  error
}

func New(model string) *Provider {
	return &Provider{model: model}
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

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(system, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		p.model,
		genai.Text(user),
		config,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
