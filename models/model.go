package models

import (
	"context"
	"errors"
)

type Model struct {
	name     string
	provider Provider
}

func New(name string, provider Provider) Model {
	return Model{
		name:     name,
		provider: provider,
	}
}

func NewModel(name string, provider Provider) Model {
	return New(name, provider)
}

func (m Model) Name() string {
	return m.name
}

func (m Model) Generate(ctx context.Context, prompt string) (string, error) {
	if m.name == "" {
		return "", errors.New("model name is required")
	}

	if m.provider == nil {
		return "", errors.New("model provider is required")
	}

	return m.provider.Generate(ctx, m.name, prompt)
}

type Provider interface {
	Generate(ctx context.Context, model string, prompt string) (string, error)
}
