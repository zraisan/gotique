package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Provider struct {
	endpoint string
	apiKey   string
}

type Option func(*Provider)

func New(apiKey string, opts ...Option) *Provider {
	provider := &Provider{
		endpoint: "https://api.openai.com/v1/responses",
		apiKey:   apiKey,
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

func WithEndpoint(endpoint string) Option {
	return func(provider *Provider) {
		provider.endpoint = endpoint
	}
}

type request struct {
	Model string `json:"model"`
	Input string `json:"input"`
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

func (p *Provider) Generate(ctx context.Context, model string, prompt string) (string, error) {
	body := request{
		Model: model,
		Input: prompt,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openai request failed: %s: %s", resp.Status, string(respBody))
	}

	var result response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.text(), nil
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
