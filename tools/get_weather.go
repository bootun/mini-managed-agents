package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bootun/mini-managed-agents/helpers"
)

func weatherAlertsTool() map[string]any {
	return helpers.ToolDefinition(
		"get_weather_alerts",
		"Get weather alerts for a US state.",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"state": map[string]any{
					"type":        "string",
					"description": "Two-letter US state code (e.g. CA, NY)",
				},
			},
			"required":             []string{"state"},
			"additionalProperties": false,
		},
	)
}

func GetWeatherAlerts(ctx context.Context, args map[string]any) (string, error) {
	state, ok := args["state"].(string)
	if !ok || strings.TrimSpace(state) == "" {
		return "", fmt.Errorf("state must be a non-empty string")
	}

	url := fmt.Sprintf("https://api.weather.gov/alerts/active/area/%s", strings.ToUpper(state))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mini-managed-agents-go/1.0")
	req.Header.Set("Accept", "application/geo+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("get weather alerts returned %s: %s", resp.Status, string(body))
	}

	randomExit()

	return string(body), nil
}
