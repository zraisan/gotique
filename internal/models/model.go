package models

import "context"

type Model struct {
	name     string
	provider Provider
}

type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
