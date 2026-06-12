// Package llm is a tiny OpenAI-compatible chat completion client. It supports
// any provider that exposes /v1/chat/completions: OpenRouter, Naraya, HF
// Router, Gemini's OpenAI-compat shim, Ollama, and self-hosted gateways.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config is the minimal set of fields the client needs.
type Config struct {
	Provider string
	Model    string
	APIKey   string
	BaseURL  string
	Lang     string // informational only
	Count    int    // informational only
	JSON     bool   // informational only
	Timeout  time.Duration
}

// Client is a single-shot LLM caller. It is safe for concurrent use.
type Client struct {
	cfg    Config
	http   *http.Client
}

// New builds a Client. Callers must set Model and APIKey on the config.
func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("llm: BaseURL is required")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("llm: APIKey is required")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("llm: Model is required")
	}
	to := cfg.Timeout
	if to <= 0 {
		to = 90 * time.Second
	}
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: to},
	}, nil
}

// Complete sends a single chat completion request and returns the model's
// first choice text. It does not stream.
func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	body := map[string]any{
		"model": c.cfg.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.8,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("llm: marshal: %w", err)
	}
	url := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("llm: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	if c.cfg.Provider == "openrouter" {
		// OpenRouter likes these for ranking/attribution
		req.Header.Set("HTTP-Referer", "https://github.com/arkydarmalik-coder/tube-trend-buddy")
		req.Header.Set("X-Title", "Tube Trend Buddy (ttb)")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: do: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("llm: HTTP %d: %s", resp.StatusCode, trimForLog(string(respBody)))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("llm: decode: %w (body=%s)", err, trimForLog(string(respBody)))
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("llm: api error: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llm: no choices in response (body=%s)", trimForLog(string(respBody)))
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func trimForLog(s string) string {
	const max = 300
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}
