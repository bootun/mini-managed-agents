package workflows

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bootun/mini-managed-agents/activities"
	"github.com/bootun/mini-managed-agents/helpers"
	"github.com/bootun/mini-managed-agents/internal/openairesp"
	"github.com/bootun/mini-managed-agents/tools"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type AgentWorkflow struct{}

func (AgentWorkflow) Run(ctx workflow.Context, input string) (string, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			// 可以为activity设定一个最大重试的次数
			MaximumAttempts: 5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	modelName := os.Getenv("OPENAI_MODEL")
	inputList := []map[string]any{{
		"type":    "message",
		"role":    "user",
		"content": input,
	}}

	for {
		var response openairesp.Response
		err := workflow.ExecuteActivity(ctx, activities.CreateOpenAIResponse, activities.OpenAIResponsesRequest{
			Model:        modelName,
			Instructions: helpers.HelpfulAgentSystemInstructions,
			Input:        inputList,
			Tools:        tools.GetTools(),
		}).Get(ctx, &response)
		if err != nil {
			return "", err
		}
		if len(response.Output) == 0 {
			if response.OutputText != "" {
				return response.OutputText, nil
			}
			return "", fmt.Errorf("openai response contained no output items")
		}

		assistantInput := appendAssistantResponse(inputList, response.Output)
		functionCalls := functionCallsFromOutput(response.Output)
		if len(functionCalls) == 0 {
			if response.OutputText != "" {
				return response.OutputText, nil
			}
			if text := firstTextFromOutput(response.Output); text != "" {
				return text, nil
			}
			return "", fmt.Errorf("openai response contained neither function calls nor final text")
		}

		inputList = assistantInput
		for _, item := range functionCalls {
			toolOutput, err := handleFunctionCall(ctx, item)
			if err != nil {
				return "", err
			}
			inputList = append(inputList, map[string]any{
				"type":    "function_call_output",
				"call_id": item.CallID,
				"output":  toolOutput,
			})
		}
	}
}

func handleFunctionCall(ctx workflow.Context, item openairesp.OutputItem) (string, error) {
	arguments := map[string]any{}
	if item.Arguments != "" {
		if err := json.Unmarshal([]byte(item.Arguments), &arguments); err != nil {
			return "", fmt.Errorf("parse tool arguments for %s: %w", item.Name, err)
		}
	}

	var result string
	err := workflow.ExecuteActivity(ctx, activities.InvokeTool, activities.ToolCallRequest{
		Name:      item.Name,
		Arguments: arguments,
	}).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

func appendAssistantResponse(inputList []map[string]any, output []openairesp.OutputItem) []map[string]any {
	updated := append([]map[string]any{}, inputList...)
	for _, item := range output {
		switch item.Type {
		case "function_call":
			updated = append(updated, map[string]any{
				"type":      item.Type,
				"call_id":   item.CallID,
				"name":      item.Name,
				"arguments": item.Arguments,
			})
		case "message":
			updated = append(updated, map[string]any{
				"type":    item.Type,
				"role":    item.Role,
				"content": item.Content,
			})
		}
	}
	return updated
}

func functionCallsFromOutput(output []openairesp.OutputItem) []openairesp.OutputItem {
	functionCalls := make([]openairesp.OutputItem, 0, len(output))
	for _, item := range output {
		if item.Type == "function_call" {
			functionCalls = append(functionCalls, item)
		}
	}
	return functionCalls
}

func firstTextFromOutput(output []openairesp.OutputItem) string {
	for _, item := range output {
		if text := firstTextFromContent(item); text != "" {
			return text
		}
	}
	return ""
}

func firstTextFromContent(item openairesp.OutputItem) string {
	for _, content := range item.Content {
		contentMap, ok := content.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := contentMap["text"].(string); ok && text != "" {
			return text
		}
		if nestedText, ok := contentMap["content"].(string); ok && nestedText != "" {
			return nestedText
		}
	}
	return ""
}
