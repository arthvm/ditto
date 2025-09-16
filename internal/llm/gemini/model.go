package gemini

import "github.com/arthvm/ditto/internal/llm"

type Model = string

const (
	GeminiFlash     Model = "gemini-2.5-flash"
	GeminiFlashLite Model = "gemini-2.5-flash-lite"
	GeminiPro       Model = "gemini-2.5-pro"
)

type provider struct {
	model Model
}

func init() {
	llm.Register("gemini", &provider{
		model: GeminiPro,
	})

	llm.Register("gemini-flash", &provider{
		model: GeminiFlash,
	})

	llm.Register("gemini-flash-lite", &provider{
		model: GeminiFlashLite,
	})
}
