package eve

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RestartService struct {
	Service        string `json:"service"`
	NamespaceAlias string `json:"namespace_alias"`
	Environment    string `json:"environment"`
}

func (r RestartService) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &r,
		validation.Field(&r.Service, validation.Required),
		validation.Field(&r.NamespaceAlias, validation.Required),
		validation.Field(&r.Environment, validation.Required),
	)
}
