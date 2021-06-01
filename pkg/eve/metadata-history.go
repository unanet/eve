package eve

import (
	"time"
)

type MetadataHistory struct {
	MetadataId  int                    `json:"metadata_id"`
	Description string                 `json:"description"`
	Value       map[string]interface{} `json:"value"`
	Created     time.Time              `json:"created"`
	CreatedBy   string                 `json:"created_by"`
	Deleted     *time.Time             `json:"deleted"`
	DeletedBy   *string                `json:"deleted_by"`
}
