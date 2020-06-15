package eve

type Environment struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Alias    string                 `json:"alias"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
