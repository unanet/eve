package eve

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Release struct {
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
	FromFeed string `json:"from_feed"`
	ToFeed   string `json:"to_feed"`
}

func (r Release) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &r,
		validation.Field(&r.Artifact, validation.Required),
		validation.Field(&r.FromFeed, validation.Required),
	)
}
