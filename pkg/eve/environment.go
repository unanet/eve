package eve

import "time"

type Environment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias,omitempty"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
