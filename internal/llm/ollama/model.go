package ollama

import "github.com/arthvm/ditto/internal/llm"

type Model = string

type provider struct {
	model Model
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
	llm.Register("ollama", &provider{
		model: GitCommitMessage,
	})
}
