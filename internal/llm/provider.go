package llm

import (
	"context"
	"errors"
	"maps"
	"slices"
)

var (
	ErrNoProvider = errors.New("no provider for name")
	ErrNoSupport  = errors.New("not supported with this model")
)

type GeneratePrParams struct {
	HeadBranch        string
	BaseBranch        string
	Log               string
	DiffStats         string
	Template          string
	Issues            []string
	AdditionalContext string
}

type Provider interface {
	GenerateCommitMessage(context.Context, string, string) (string, error)
	GeneratePr(context.Context, GeneratePrParams) (string, error)
}

var providers map[string]Provider

func init() {
	providers = map[string]Provider{}
}

func Register(name string, provider Provider) {
	providers[name] = provider
}

func ListProviders() []string {
	return slices.Collect(maps.Keys(providers))
}

func GetProvider(name string) (Provider, error) {
	provider, exists := providers[name]
	if !exists {
		return nil, ErrNoProvider
	}

	return provider, nil
}
