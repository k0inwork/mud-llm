package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendPromptWithCustomEnv(t *testing.T) {
	// This test uses a mock server to verify that the client correctly uses custom environment variables.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedAuthHeader := "Bearer test-key"
		if r.Header.Get("Authorization") != expectedAuthHeader {
			t.Errorf("Expected Authorization header %s, got %s", expectedAuthHeader, r.Header.Get("Authorization"))
		}
		var reqBody LLMRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody.Model != "test-model" {
			t.Errorf("Expected model test-model, got %s", reqBody.Model)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LLMResponse{
			Choices: []Choice{
				{
					Message: Message{
						Content: `{"narrative": "Custom env test narrative.", "tool_calls": []}`,
					},
				},
			},
		})
	}))
	defer mockServer.Close()

	t.Setenv("LLM_API_ENDPOINT", mockServer.URL)
	t.Setenv("LLM_API_KEY", "test-key")
	t.Setenv("LLM_MODEL_NAME", "test-model")

	client := NewClient()
	resp, err := client.SendPrompt(context.Background(), "test prompt")
	if err != nil {
		t.Fatalf("SendPrompt failed: %v", err)
	}
	if resp.Narrative != "Custom env test narrative." {
		t.Errorf("Expected narrative 'Custom env test narrative.', got '%s'", resp.Narrative)
	}
}

func TestSendPromptWithDefaultValues_Mocked(t *testing.T) {
	// This test uses a mock server to verify that the client correctly uses default values when environment variables are not set.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedAuthHeader := "Bearer unused"
		if r.Header.Get("Authorization") != expectedAuthHeader {
			t.Errorf("Expected Authorization header %s, got %s", expectedAuthHeader, r.Header.Get("Authorization"))
		}
		var reqBody LLMRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		expectedModel := "gpt-4.1-2025-04-14"
		if reqBody.Model != expectedModel {
			t.Errorf("Expected model %s, got %s", expectedModel, reqBody.Model)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LLMResponse{
			Choices: []Choice{
				{
					Message: Message{
						Content: `{"narrative": "Default values mock test narrative.", "tool_calls": []}`,
					},
				},
			},
		})
	}))
	defer mockServer.Close()

	// Temporarily set the endpoint to our mock server, but leave others unset to test defaults.
	originalEndpoint := os.Getenv("LLM_API_ENDPOINT")
	os.Setenv("LLM_API_ENDPOINT", mockServer.URL)
	defer os.Setenv("LLM_API_ENDPOINT", originalEndpoint)

	originalApiKey := os.Getenv("LLM_API_KEY")
	os.Unsetenv("LLM_API_KEY")
	defer os.Setenv("LLM_API_KEY", originalApiKey)

	originalModelName := os.Getenv("LLM_MODEL_NAME")
	os.Unsetenv("LLM_MODEL_NAME")
	defer os.Setenv("LLM_MODEL_NAME", originalModelName)

	client := NewClient()
	resp, err := client.SendPrompt(context.Background(), "test prompt for defaults")
	if err != nil {
		t.Fatalf("SendPrompt with defaults failed: %v", err)
	}
	if resp.Narrative != "Default values mock test narrative." {
		t.Errorf("Expected narrative 'Default values mock test narrative.', got '%s'", resp.Narrative)
	}
}


