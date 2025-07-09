//go:build integration

package llm

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSendPrompt_Integration(t *testing.T) {
	// This test will be skipped if the integration tag is not provided.
	// To run this test, use: go test -tags=integration ./...

	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" || apiKey == "unused" {
		t.Skip("Skipping integration test: LLM_API_KEY is not set or is 'unused'")
	}

	// Use the default values by not setting the environment variables
	// The NewClient function will apply the defaults
	client := NewClient()

	prompt := "This is a test prompt to the real LLM. Please respond with a short narrative and no tool calls."

	resp, err := client.SendPrompt(context.Background(), prompt)

	if err != nil {
		t.Fatalf("SendPrompt failed with error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected a response from the LLM, but got nil")
	}

	if resp.Narrative == "" {
		t.Error("Expected a non-empty narrative in the response")
	}

	t.Logf("Received narrative: %s", resp.Narrative)
}

func TestSendPrompt_RealNetwork_DefaultValues(t *testing.T) {
	// This test makes a real network request to the default LLM service.
	// It will fail if there is no network access or the service is down.
	// It is included to run always, per user request.
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	// Ensure environment variables are unset so defaults are used.
	originalEndpoint := os.Getenv("LLM_API_ENDPOINT")
	os.Unsetenv("LLM_API_ENDPOINT")
	defer os.Setenv("LLM_API_ENDPOINT", originalEndpoint)

	originalApiKey := os.Getenv("LLM_API_KEY")
	os.Unsetenv("LLM_API_KEY")
	defer os.Setenv("LLM_API_KEY", originalApiKey)

	originalModelName := os.Getenv("LLM_MODEL_NAME")
	os.Unsetenv("LLM_MODEL_NAME")
	defer os.Setenv("LLM_MODEL_NAME", originalModelName)

	client := NewClient()

	// Verify that the client has loaded the default values
	if client.apiKey != "unused" {
		t.Fatalf("Expected default apiKey 'unused', but got '%s'", client.apiKey)
	}
	if client.apiURL != "https://api.llm7.io/v1" {
		t.Fatalf("Expected default apiURL 'https://api.llm7.io/v1', but got '%s'", client.apiURL)
	}

	prompt := "This is a real network test from a Go test suite. Please respond with a short, simple confirmation message. No tool calls are needed."

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.SendPrompt(ctx, prompt)

	if err != nil {
		t.Fatalf("SendPrompt to real network failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected a response from the real LLM, but got nil")
	}

	if resp.Narrative == "" {
		t.Error("Expected a non-empty narrative in the response from the real LLM")
	}

	t.Logf("Successfully received narrative from real LLM: %s", resp.Narrative)
}