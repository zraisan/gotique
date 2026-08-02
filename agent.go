package gotique

import (
	"context"

	"github.com/google/uuid"
	"github.com/zraisan/gotique/models"
)

type Agent struct {
	Name               string
	Instructions       []string
	SystemPrompt       string
	Model              models.Model
	FallbackModels     []models.Model
	NumHistoryMessages int
	Sessions           map[uuid.UUID]*Session
}

type Session struct {
	SessionID uuid.UUID
	Messages  []Message
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Prompt struct {
	Role    Role
	Content string
}

type Message struct {
	prompts [2]Prompt
}

func NewSession() uuid.UUID {
	return uuid.New()
}

func NewAgent(a Agent) *Agent {
	var session Session
	session.SessionID = NewSession()
	a.Sessions = make(map[uuid.UUID]*Session)
	a.Sessions[session.SessionID] = &session
	return &a
}

func (a *Agent) Run(ctx context.Context, prompt string, sessionID uuid.UUID) (string, error) {
	resp, err := a.Model.Generate(ctx, models.Request{
		SystemPrompt: a.SystemPrompt,
		Instructions: a.Instructions,
		Input:        prompt,
	})
	if err != nil {
		for i := range a.FallbackModels {
			resp, err = a.FallbackModels[i].Generate(ctx, models.Request{
				SystemPrompt: a.SystemPrompt,
				Instructions: a.Instructions,
				Input:        prompt,
			})

			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return "", err
	}

	if a.NumHistoryMessages > 0 {
		message := Message{
			prompts: [2]Prompt{
				{Role: RoleUser, Content: prompt},
				{Role: RoleAssistant, Content: resp.Text},
			},
		}
		a.Sessions[sessionID].Messages = append(a.Sessions[sessionID].Messages, message)
	}

	return resp.Text, nil
}
