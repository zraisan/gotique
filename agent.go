package gotique

import (
	"context"

	"github.com/zraisan/gotique/models"
)

type Agent struct {
	name         string
	instructions []string
	systemPrompt string
	model        models.Model
}

type AgentOption func(*Agent)

func NewAgent(opts ...AgentOption) *Agent {
	agent := &Agent{}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func WithName(name string) AgentOption {
	return func(agent *Agent) {
		agent.name = name
	}
}

func WithInstructions(instructions ...string) AgentOption {
	return func(agent *Agent) {
		agent.instructions = instructions
	}
}

func WithSystemPrompt(systemPrompt string) AgentOption {
	return func(agent *Agent) {
		agent.systemPrompt = systemPrompt
	}
}

func WithModel(model models.Model) AgentOption {
	return func(agent *Agent) {
		agent.model = model
	}
}

func (a *Agent) SetModel(model models.Model) {
	a.model = model
}

func (a *Agent) NewModel(name string, provider models.Provider) error {
	a.model = models.New(name, provider)
	return nil
}

func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	return a.model.Generate(ctx, prompt)
}
