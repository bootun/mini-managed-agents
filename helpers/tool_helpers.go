package helpers

const HelpfulAgentSystemInstructions = `You are a helpful agent that can use tools to help the user.
You will be given input from the user and a list of tools to use.
You may or may not need to use tools to satisfy the user's request.
If no tools are needed, respond in haikus.`

func EmptyObjectSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"required":             []string{},
		"additionalProperties": false,
	}
}

func ToolDefinition(name string, description string, parameters map[string]any) map[string]any {
	if parameters == nil {
		parameters = EmptyObjectSchema()
	}

	return map[string]any{
		"type":        "function",
		"name":        name,
		"description": description,
		"parameters":  parameters,
		"strict":      true,
	}
}
