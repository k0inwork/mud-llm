package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = "unused"
	}
	apiURL := os.Getenv("LLM_API_ENDPOINT")
	if apiURL == "" {
		apiURL = "https://api.llm7.io/v1"
	}

	return &Client{
		apiKey: apiKey,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type LLMRequest struct {
	Model         string        `json:"model"`
	Messages      []Message     `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type InnerLLMResponse struct {
	Narrative string      `json:"narrative"`
	ToolCalls []ToolCall  `json:"tool_calls"`
}

type ToolCall struct {
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (c *Client) SendPrompt(ctx context.Context, prompt string) (*InnerLLMResponse, error) {
	modelName := os.Getenv("LLM_MODEL_NAME")
	if modelName == "" {
		modelName = "gpt-4.1-2025-04-14"
	}

	reqBody := LLMRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant for a multi-user dungeon game. Your responses should be in JSON format, with a 'narrative' field for text to be shown to the player, and a 'tool_calls' field for any actions the AI should take.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		ResponseFormat: &ResponseFormat{Type: "json_object"},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	requestURL := c.apiURL + "/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the raw response body for debugging
	fmt.Printf("Raw LLM Response: %s\n", string(bodyBytes))

	var llmResponse LLMResponse
	if err := json.Unmarshal(bodyBytes, &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	if len(llmResponse.Choices) == 0 {
		return nil, errors.New("no choices in LLM response")
	}

	content := llmResponse.Choices[0].Message.Content

	var innerLLMResponse InnerLLMResponse
	if err := json.Unmarshal([]byte(content), &innerLLMResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner LLM response: %w", err)
	}

	return &innerLLMResponse, nil
}