package copilot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultClientID = "Iv23liTZx6XX7PidTYbw"
	deviceCodeURL   = "https://github.com/login/device/code"
	tokenURL        = "https://github.com/login/oauth/access_token"
	verificationURL = "https://github.com/login/device"
	deviceGrantType = "urn:ietf:params:oauth:grant-type:device_code"
	tokenFileName   = "github_token"
)

type deviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

// resolveGitHubToken returns a GitHub token suitable for the Copilot
// API. It checks, in order:
//  1. Explicit API key from config (copilot.api_key)
//  2. Previously saved ditto token (~/.config/ditto/github_token)
//  3. OAuth device flow (interactive, saves token for next time)
func resolveGitHubToken(apiKey, clientID string) (string, error) {
	if apiKey != "" {
		return apiKey, nil
	}

	if token, err := loadSavedToken(); err == nil && token != "" {
		return token, nil
	}

	if clientID == "" {
		clientID = defaultClientID
	}
	token, err := deviceFlowAuth(clientID)
	if err != nil {
		return "", fmt.Errorf("copilot auth: %w", err)
	}

	_ = saveToken(token)

	return token, nil
}

func loadSavedToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, ".config", "ditto", tokenFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	if token := strings.TrimSpace(string(data)); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("empty token file")
}

func saveToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".config", "ditto")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, tokenFileName), []byte(token), 0o600)
}

// deviceFlowAuth implements the OAuth 2.0 Device Authorization Grant
// (RFC 8628) against GitHub. It prints a user code, attempts to open
// the verification URL in a browser, and polls until the user authorises
// or the code expires.
func deviceFlowAuth(clientID string) (string, error) {
	code, err := requestDeviceCode(clientID)
	if err != nil {
		return "", err
	}

	fmt.Printf("\nTo authenticate with GitHub Copilot:\n")
	fmt.Printf("  1. Open:  %s\n", verificationURL)
	fmt.Printf("  2. Enter: %s\n\n", code.UserCode)

	_ = openBrowser(verificationURL)

	return pollForToken(clientID, code)
}

func requestDeviceCode(clientID string) (*deviceCodeResponse, error) {
	form := url.Values{
		"client_id": {clientID},
		"scope":     {"read:user"},
	}

	req, err := http.NewRequest(http.MethodPost, deviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create device code request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device code request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device code request: status %d", resp.StatusCode)
	}

	var code deviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&code); err != nil {
		return nil, fmt.Errorf("decode device code response: %w", err)
	}

	if code.Interval == 0 {
		code.Interval = 5
	}

	return &code, nil
}

func pollForToken(clientID string, code *deviceCodeResponse) (string, error) {
	interval := time.Duration(code.Interval) * time.Second
	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("device code expired, please try again")
		}

		time.Sleep(interval)

		token, err := exchangeDeviceCode(clientID, code.DeviceCode)
		if err != nil {
			return "", err
		}

		switch token.Error {
		case "":
			return token.AccessToken, nil
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
		case "expired_token":
			return "", fmt.Errorf("device code expired, please try again")
		case "access_denied":
			return "", fmt.Errorf("access denied by user")
		default:
			return "", fmt.Errorf("unexpected error: %s", token.Error)
		}
	}
}

func exchangeDeviceCode(clientID, deviceCode string) (*tokenResponse, error) {
	form := url.Values{
		"client_id":   {clientID},
		"device_code": {deviceCode},
		"grant_type":  {deviceGrantType},
	}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
