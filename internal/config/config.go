package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider   string       `yaml:"provider"`
	Model      string       `yaml:"model"`
	BaseBranch string       `yaml:"base_branch"`
	LLM        LLMConfig    `yaml:"llm"`
	Commit     CommitConfig `yaml:"commit"`
	PR         PRConfig     `yaml:"pr"`
	Gemini     GeminiConfig `yaml:"gemini"`
	Ollama     OllamaConfig `yaml:"ollama"`
}

type LLMConfig struct {
	Timeout     time.Duration `yaml:"timeout"`
	Temperature float32       `yaml:"temperature"`
}

func (l *LLMConfig) UnmarshalYAML(value *yaml.Node) error {
	type plain struct {
		Timeout     string  `yaml:"timeout"`
		Temperature float32 `yaml:"temperature"`
	}

	var p plain
	if err := value.Decode(&p); err != nil {
		return err
	}

	l.Temperature = p.Temperature
	if p.Timeout != "" {
		d, err := time.ParseDuration(p.Timeout)
		if err != nil {
			return fmt.Errorf("llm.timeout: %w", err)
		}
		l.Timeout = d
	}

	return nil
}

type CommitConfig struct {
	Prompt string `yaml:"prompt"`
}

type PRConfig struct {
	Prompt       string `yaml:"prompt"`
	TemplatePath string `yaml:"template_path"`
}

type GeminiConfig struct {
	APIKey string `yaml:"api_key"`
}

type OllamaConfig struct {
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

func defaults() Config {
	return Config{
		Provider:   "gemini",
		Model:      "gemini-2.5-flash",
		BaseBranch: "main",
		LLM: LLMConfig{
			Timeout: 2 * time.Minute,
		},
		Ollama: OllamaConfig{
			Host:  "http://localhost:11434",
			Model: "tavernari/git-commit-message",
		},
	}
}

func Load(repoRoot string) (Config, error) {
	cfg := defaults()

	userDir, err := os.UserConfigDir()
	if err == nil {
		userPath := filepath.Join(userDir, "ditto", "config.yaml")
		if err := mergeFromFile(&cfg, userPath); err != nil {
			return cfg, fmt.Errorf("user config: %w", err)
		}
	}

	if repoRoot != "" {
		projectPath := filepath.Join(repoRoot, ".ditto.yaml")
		if err := mergeFromFile(&cfg, projectPath); err != nil {
			return cfg, fmt.Errorf("project config: %w", err)
		}
	}

	mergeFromEnv(&cfg)

	return cfg, nil
}

func mergeFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func mergeFromEnv(cfg *Config) {
	if v, ok := os.LookupEnv("DITTO_PROVIDER"); ok {
		cfg.Provider = v
	}
	if v, ok := os.LookupEnv("DITTO_MODEL"); ok {
		cfg.Model = v
	}
	if v, ok := os.LookupEnv("DITTO_BASE_BRANCH"); ok {
		cfg.BaseBranch = v
	}
	if v, ok := os.LookupEnv("DITTO_LLM_TIMEOUT"); ok {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.LLM.Timeout = d
		}
	}
	if v, ok := os.LookupEnv("DITTO_LLM_TEMPERATURE"); ok {
		var t float32
		if _, err := fmt.Sscanf(v, "%g", &t); err == nil {
			cfg.LLM.Temperature = t
		}
	}
	if v, ok := os.LookupEnv("GEMINI_API_KEY"); ok {
		cfg.Gemini.APIKey = v
	}
	if v, ok := os.LookupEnv("OLLAMA_HOST"); ok {
		cfg.Ollama.Host = v
	}
	if v, ok := os.LookupEnv("OLLAMA_MODEL"); ok {
		cfg.Ollama.Model = v
	}
}
