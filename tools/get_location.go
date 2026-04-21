package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bootun/mini-managed-agents/helpers"
)

func getLocationTool() map[string]any {
	return helpers.ToolDefinition(
		"get_location_info",
		"Get the location information for an IP address. This includes the city, state, and country.",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"ipaddress": map[string]any{
					"type":        "string",
					"description": "An IP address",
				},
			},
			"required":             []string{"ipaddress"},
			"additionalProperties": false,
		},
	)
}

func getIPAddressTool() map[string]any {
	return helpers.ToolDefinition(
		"get_ip_address",
		"Get the IP address of the current machine.",
		nil,
	)
}

func GetIPAddress(ctx context.Context, args map[string]any) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://icanhazip.com", nil)
	if err != nil {
		return "", err
	}

	randomExit()

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
		return "", fmt.Errorf("get ip address returned %s: %s", resp.Status, string(body))
	}

	return strings.TrimSpace(string(body)), nil
}

func GetLocationInfo(ctx context.Context, args map[string]any) (string, error) {
	ipAddress, ok := args["ipaddress"].(string)
	if !ok || strings.TrimSpace(ipAddress) == "" {
		return "", fmt.Errorf("ipaddress must be a non-empty string")
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	randomExit()

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
		return "", fmt.Errorf("get location returned %s: %s", resp.Status, string(body))
	}

	var payload struct {
		City       string `json:"city"`
		RegionName string `json:"regionName"`
		Country    string `json:"country"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s, %s", payload.City, payload.RegionName, payload.Country), nil
}
