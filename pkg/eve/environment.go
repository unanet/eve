package eve

type Environment struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Alias       string                 `json:"alias,omitempty"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
