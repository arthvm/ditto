package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/arthvm/ditto/internal/llm"
)

type Model = string

type model struct {
}

const (
	GitCommitMessage Model = "tavernari/git-commit-message"
)

type generateRequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Raw    bool   `json:"raw"`
}

type generateResponseBody struct {
	Response string `json:"response"`
}

func init() {
	llm.Register("ollama", &model{})
}

func (m *model) GenerateCommitMessage(
	ctx context.Context,
	diff string,
	additionalContext string,
) (string, error) {
	baseUrl, exists := os.LookupEnv("OLLAMA_HOST")
	if !exists {
		baseUrl = "http://localhost:11434"
	}

	url := fmt.Sprintf("%s/api/generate", baseUrl)

	body := generateRequestBody{
		Model:  GitCommitMessage,
		Prompt: diff,
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
