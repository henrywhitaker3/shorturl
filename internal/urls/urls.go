package urls

import (
	"context"

	"github.com/henrywhitaker3/shorturl/internal/uuid"
)

type Urls interface {
	Create(context.Context, CreateParams) (*Url, error)
	Get(context.Context, uuid.UUID) (*Url, error)
	GetAlias(context.Context, string) (*Url, error)
}
