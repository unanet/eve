package eve

import (
	"time"
)

type Cluster struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ProviderGroup string    `json:"provider_group"`
	SchQueueUrl   string    `json:"sch_queue_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
