package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

func (p *provider) Generate(ctx context.Context, system, user string) (string, error) {
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
