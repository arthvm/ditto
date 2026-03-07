package copilot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.githubcopilot.com"

// Provider implements workflow.Provider using the GitHub Copilot API.
// It authenticates with a GitHub OAuth token obtained via the device flow,
// stored in the system keychain. If the token becomes invalid, it
// re-authenticates automatically without user intervention.
type Provider struct {
	model       string
	temperature float32
	token       string
	clientID    string
}

func New(model string, temperature float32, apiKey, clientID string) (*Provider, error) {
	token, err := resolveGitHubToken(apiKey, clientID)
	if err != nil {
		return nil, err
	}

	return &Provider{
		model:       model,
		temperature: temperature,
		token:       token,
		clientID:    clientID,
	}, nil
}

func (p *Provider) Generate(ctx context.Context, system, user string) (string, error) {
	result, err := p.generate(ctx, system, user)
	if err == nil {
		return result, nil
	}

	// On auth failure, clear the stored token and re-authenticate once.
	if isAuthError(err) {
		clearStoredToken()
		newToken, authErr := runDeviceFlow(p.clientID)
		if authErr != nil {
			return "", fmt.Errorf("re-authentication failed: %w", authErr)
		}
		p.token = newToken
		return p.generate(ctx, system, user)
	}

	return "", err
}

func (p *Provider) generate(ctx context.Context, system, user string) (string, error) {
	body, err := p.buildRequest(system, user)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", body)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	p.setHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("copilot request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", &authError{status: resp.StatusCode}
	}

	if resp.StatusCode != http.StatusOK {
		// Read error is intentionally ignored: the status code is already
		// informative and a body read failure would obscure the real error.
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("copilot api: status %d: %s", resp.StatusCode, errBody)
	}

	var chatResp chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("copilot api: empty response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

type authError struct {
	status int
}

func (e *authError) Error() string {
	return fmt.Sprintf("copilot api: auth error (status %d)", e.status)
}

func isAuthError(err error) bool {
	return errors.As(err, new(*authError))
}

func (p *Provider) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Editor-Version", "Ditto/1.0")
	req.Header.Set("Editor-Plugin-Version", "Ditto/1.0")
	req.Header.Set("Copilot-Integration-Id", "vscode-chat")
}

func (p *Provider) buildRequest(system, user string) (*bytes.Buffer, error) {
	messages := []chatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	}

	reqBody := chatCompletionRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	}
	if p.temperature != 0 {
		reqBody.Temperature = &p.temperature
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return nil, err
	}
	return buf, nil
}

// OpenAI-compatible request/response types.

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Temperature *float32      `json:"temperature,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
