package openairesp

type Response struct {
	Output     []OutputItem `json:"output"`
	OutputText string       `json:"output_text"`
}

type OutputItem struct {
	Type      string `json:"type"`
	ID        string `json:"id,omitempty"`
	CallID    string `json:"call_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
	Role      string `json:"role,omitempty"`
	Status    string `json:"status,omitempty"`
	Content   []any  `json:"content,omitempty"`
}

type Request struct {
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Input        []map[string]any `json:"input"`
	Tools        []map[string]any `json:"tools,omitempty"`
}
