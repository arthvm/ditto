package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Provider implements workflow.Provider using the Ollama local API.
type Provider struct {
	host        string
	model       string
	temperature float32
}

// New creates an Ollama provider. A temperature of 0 uses the model's default.
func New(host, model string, temperature float32) *Provider {
	return &Provider{host: host, model: model, temperature: temperature}
}

type generateOptions struct {
	Temperature float32 `json:"temperature,omitempty"`
}

type generateRequestBody struct {
	Model   string          `json:"model"`
	System  string          `json:"system"`
	Prompt  string          `json:"prompt"`
	Stream  bool            `json:"stream"`
	Raw     bool            `json:"raw"`
	Options generateOptions `json:"options,omitempty"`
}

type generateResponseBody struct {
	Response string `json:"response"`
}

func (p *Provider) Generate(ctx context.Context, system, user string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", p.host)

	body := generateRequestBody{
		Model:   p.model,
		System:  system,
		Prompt:  user,
		Stream:  false,
		Raw:     false,
		Options: generateOptions{Temperature: p.temperature},
	}
	bodyBuf := &bytes.Buffer{}

	if err := json.NewEncoder(bodyBuf).Encode(body); err != nil {
		return "", fmt.Errorf("encode body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyBuf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	var resBody generateResponseBody

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", fmt.Errorf("decode body: %w", err)
	}

	return resBody.Response, nil
}
