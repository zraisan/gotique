package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zraisan/gotique/models"
)

type Provider struct {
	endpoint string
	apiKey   string
}

type request struct {
	Model        string `json:"model"`
	Instructions string `json:"instructions,omitempty"`
	Input        string `json:"input"`
}

type response struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func New(apiKey string, endpoint ...string) *Provider {
	providerEndpoint := "https://api.openai.com/v1/responses"
	if len(endpoint) > 0 && endpoint[0] != "" {
		providerEndpoint = endpoint[0]
	}

	return &Provider{
		endpoint: providerEndpoint,
		apiKey:   apiKey,
	}
}

func (p *Provider) Generate(ctx context.Context, modelReq models.Request) (*models.Response, error) {
	instructionParts := []string{}
	if modelReq.SystemPrompt != "" {
		instructionParts = append(instructionParts, modelReq.SystemPrompt)
	}

	instructionParts = append(instructionParts, modelReq.Instructions...)

	body := request{
		Model:        modelReq.Model,
		Instructions: strings.Join(instructionParts, "\n"),
		Input:        modelReq.Input,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("openai request failed: %s: %s", resp.Status, string(respBody))
	}

	var result response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &models.Response{
		Text: result.text(),
	}, nil
}

func (r *response) text() string {
	if r.OutputText != "" {
		return r.OutputText
	}

	var parts []string
	for _, output := range r.Output {
		if output.Type != "message" {
			continue
		}

		for _, content := range output.Content {
			if content.Type == "output_text" {
				parts = append(parts, content.Text)
			}
		}
	}

	return strings.Join(parts, "")
}
