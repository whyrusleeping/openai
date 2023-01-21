package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Auth string
	Org  string
}

type CompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type CompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []CompChoice `json:"choices"`
	Usage   *Usage       `json:"usage"`
}

type CompChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ApiError struct {
	ErrorVal struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (ae ApiError) Error() string {
	return fmt.Sprintf("ApiError(%s): %s", ae.ErrorVal.Type, ae.ErrorVal.Message)
}

func (c *Client) Completion(ctx context.Context, creq *CompletionRequest) (*CompletionResponse, error) {
	b, err := json.Marshal(creq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth)
	req.Header.Set("OpenAI-Organization", c.Org)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var out ApiError
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("invalid error body (status %d): %w", resp.StatusCode, err)
		}
		return nil, out
	}

	var out CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
