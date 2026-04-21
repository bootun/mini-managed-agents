package activities

import (
	"context"
	"fmt"

	"github.com/bootun/mini-managed-agents/tools"
)

type ToolCallRequest struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

func InvokeTool(ctx context.Context, request ToolCallRequest) (string, error) {
	handler, err := tools.GetHandler(request.Name)
	if err != nil {
		return "", err
	}

	result, err := handler(ctx, request.Arguments)
	if err != nil {
		return "", fmt.Errorf("invoke tool %s: %w", request.Name, err)
	}

	return result, nil
}
