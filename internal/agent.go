package internal

import "github.com/zraisan/gotique/internal/models"

type Agent struct {
	name         string
	instructions []string
	systemPrompt string
}

func (a *Agent) NewModel(name string, provider models.Provider) error {
	return nil
}
