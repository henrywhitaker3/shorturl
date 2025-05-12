package urls

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
)

type Service struct {
	db *queries.Queries
}

type ServiceOpts struct {
	DB *queries.Queries
}

func New(opts ServiceOpts) *Service {
	return &Service{
		db: opts.DB,
	}
}

type CreateParams struct {
	ID     uuid.UUID
	Url    string
	Domain string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Url, error) {
	url, err := s.db.CreateUrl(ctx, queries.CreateUrlParams{
		ID:     params.ID.UUID(),
		Url:    params.Url,
		Domain: params.Domain,
	})
	if err != nil {
		return nil, fmt.Errorf("store url: %w", err)
	}

	return mapUrl(url), nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Url, error) {
	url, err := s.db.GetUrl(ctx, id.UUID())
	if err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	}
	return mapUrl(url), nil
}

func (s *Service) GetAlias(ctx context.Context, alias string) (*Url, error) {
	url, err := s.db.GetUrlByAlias(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("get url by alias: %w", err)
	}
	return mapUrl(url), nil
}

var _ Urls = &Service{}
