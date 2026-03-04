package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (p *provider) Generate(ctx context.Context, system, user string) (string, error) {
	baseURL, exists := os.LookupEnv("OLLAMA_HOST")
	if !exists {
		baseURL = "http://localhost:11434"
	}

	url := fmt.Sprintf("%s/api/generate", baseURL)

	body := generateRequestBody{
		Model:  p.model,
		System: system,
		Prompt: user,
		Stream: false,
		Raw:    false,
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
