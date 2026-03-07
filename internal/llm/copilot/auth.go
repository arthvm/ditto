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

	"github.com/zalando/go-keyring"
)

// authClient is used for all OAuth authentication requests. A fixed timeout
// prevents these calls from hanging indefinitely on network issues.
var authClient = &http.Client{Timeout: 15 * time.Second}

// warnf prints a non-fatal warning to stderr.
func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

const (
	defaultClientID  = "Iv23liTZx6XX7PidTYbw"
	deviceCodeURL    = "https://github.com/login/device/code"
	tokenURL         = "https://github.com/login/oauth/access_token"
	verificationURL  = "https://github.com/login/device"
	deviceGrantType  = "urn:ietf:params:oauth:grant-type:device_code"
	refreshGrantType = "refresh_token"
	keychainService  = "ditto"
	keychainAccount  = "github_oauth_token"
	legacyTokenFile  = "github_token"
)

// storedToken is persisted as JSON in the system keychain.
type storedToken struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	ExpiresAt             time.Time `json:"expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (t *storedToken) accessTokenValid() bool {
	return t.ExpiresAt.IsZero() || time.Now().Before(t.ExpiresAt.Add(-time.Minute))
}

func (t *storedToken) refreshTokenValid() bool {
	return t.RefreshToken != "" &&
		(t.RefreshTokenExpiresAt.IsZero() || time.Now().Before(t.RefreshTokenExpiresAt))
}

type deviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Error                 string `json:"error"`
}

// resolveGitHubToken returns a valid GitHub OAuth token for the Copilot API.
// It checks, in order:
//  1. Explicit API key from config (copilot.api_key)
//  2. Stored token from keychain — refreshes it silently if expired
//  3. Legacy plain-text file (~/.config/ditto/github_token) — migrates on find
//  4. OAuth device flow (interactive, one-time setup)
func resolveGitHubToken(apiKey, clientID string) (string, error) {
	if apiKey != "" {
		return apiKey, nil
	}

	if clientID == "" {
		clientID = defaultClientID
	}

	if stored, err := loadStoredToken(); err == nil {
		if stored.accessTokenValid() {
			return stored.AccessToken, nil
		}
		if stored.refreshTokenValid() {
			refreshed, err := refreshAccessToken(clientID, stored.RefreshToken)
			if err == nil {
				if err := saveStoredToken(refreshed); err != nil {
					warnf("could not save refreshed token to keychain: %v", err)
				}
				return refreshed.AccessToken, nil
			}
			// Refresh failed — fall through to device flow.
		}
	}

	// Migrate from legacy plain-text file if it exists.
	if token, err := loadLegacyToken(); err == nil && token != "" {
		stored := &storedToken{AccessToken: token}
		if err := saveStoredToken(stored); err != nil {
			warnf("could not migrate token to keychain: %v", err)
		}
		if err := removeLegacyToken(); err != nil {
			warnf("could not remove legacy token file: %v", err)
		}
		return token, nil
	}

	return runDeviceFlow(clientID)
}

// runDeviceFlow runs the OAuth device flow and saves the resulting token to the keychain.
func runDeviceFlow(clientID string) (string, error) {
	token, err := deviceFlowAuth(clientID)
	if err != nil {
		return "", fmt.Errorf("copilot auth: %w", err)
	}
	if err := saveStoredToken(token); err != nil {
		warnf("could not save token to keychain: %v", err)
	}
	return token.AccessToken, nil
}

// clearStoredToken removes the token from the keychain so the next call
// to resolveGitHubToken triggers a fresh device flow.
func clearStoredToken() {
	if err := keyring.Delete(keychainService, keychainAccount); err != nil {
		warnf("could not clear stored token from keychain: %v", err)
	}
}

func loadStoredToken() (*storedToken, error) {
	raw, err := keyring.Get(keychainService, keychainAccount)
	if err != nil {
		return nil, err
	}
	var stored storedToken
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		// Legacy: plain string token (no expiry info).
		if strings.TrimSpace(raw) != "" {
			return &storedToken{AccessToken: strings.TrimSpace(raw)}, nil
		}
		return nil, err
	}
	if stored.AccessToken == "" {
		return nil, fmt.Errorf("keychain entry has empty access token")
	}
	return &stored, nil
}

func saveStoredToken(t *storedToken) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return keyring.Set(keychainService, keychainAccount, string(data))
}

func refreshAccessToken(clientID, refreshToken string) (*storedToken, error) {
	form := url.Values{
		"client_id":     {clientID},
		"grant_type":    {refreshGrantType},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := authClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh request: status %d", resp.StatusCode)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode refresh response: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("refresh error: %s", tokenResp.Error)
	}

	return tokenResponseToStored(&tokenResp), nil
}

func loadLegacyToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(home, ".config", "ditto", legacyTokenFile))
	if err != nil {
		return "", err
	}
	if token := strings.TrimSpace(string(data)); token != "" {
		return token, nil
	}
	return "", fmt.Errorf("empty token file")
}

func removeLegacyToken() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(home, ".config", "ditto", legacyTokenFile))
}

// deviceFlowAuth implements the OAuth 2.0 Device Authorization Grant (RFC 8628).
func deviceFlowAuth(clientID string) (*storedToken, error) {
	code, err := requestDeviceCode(clientID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nTo authenticate with GitHub Copilot:\n")
	fmt.Printf("  1. Open:  %s\n", verificationURL)
	fmt.Printf("  2. Enter: %s\n\n", code.UserCode)

	_ = openBrowser(verificationURL) // best-effort; URL is already printed above

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

	resp, err := authClient.Do(req)
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

func pollForToken(clientID string, code *deviceCodeResponse) (*storedToken, error) {
	interval := time.Duration(code.Interval) * time.Second
	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("device code expired, please try again")
		}

		time.Sleep(interval)

		tokenResp, err := exchangeDeviceCode(clientID, code.DeviceCode)
		if err != nil {
			return nil, err
		}

		switch tokenResp.Error {
		case "":
			return tokenResponseToStored(tokenResp), nil
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
		case "expired_token":
			return nil, fmt.Errorf("device code expired, please try again")
		case "access_denied":
			return nil, fmt.Errorf("access denied by user")
		default:
			return nil, fmt.Errorf("unexpected error: %s", tokenResp.Error)
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

	resp, err := authClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request: status %d", resp.StatusCode)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

func tokenResponseToStored(r *tokenResponse) *storedToken {
	s := &storedToken{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}
	now := time.Now()
	if r.ExpiresIn > 0 {
		s.ExpiresAt = now.Add(time.Duration(r.ExpiresIn) * time.Second)
	}
	if r.RefreshTokenExpiresIn > 0 {
		s.RefreshTokenExpiresAt = now.Add(time.Duration(r.RefreshTokenExpiresIn) * time.Second)
	}
	return s
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
