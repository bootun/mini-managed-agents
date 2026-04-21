package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bootun/mini-managed-agents/internal/openairesp"
)

type OpenAIResponsesRequest struct {
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Input        []map[string]any `json:"input"`
	Tools        []map[string]any `json:"tools"`
}

// 请求OpenAI Responses API生成响应
func CreateOpenAIResponse(ctx context.Context, request OpenAIResponsesRequest) (*openairesp.Response, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")

	body, err := json.Marshal(openairesp.Request{
		Model:        request.Model,
		Instructions: request.Instructions,
		Input:        request.Input,
		Tools:        request.Tools,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal openai request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/responses", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create openai request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("call openai responses api: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read openai response: %w", err)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("openai responses api returned %s: %s", response.Status, string(responseBody))
	}

	rawBody := string(responseBody)
	if len(rawBody) > 2000 {
		rawBody = rawBody[:2000]
	}

	var parsed openairesp.Response
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf(
			"unmarshal openai response: status=%s content-type=%q body=%q err=%w",
			response.Status,
			response.Header.Get("Content-Type"),
			rawBody,
			err,
		)
	}

	log.Printf("oai request: %v", string(body))
	log.Printf("oai response: %v", parsed)

	return &parsed, nil
}
