package urls

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
)

type Service struct {
	db    *queries.Queries
	conn  *sql.DB
	alias *Alias
}

type ServiceOpts struct {
	DB    *queries.Queries
	Conn  *sql.DB
	Alias *Alias
}

func New(opts ServiceOpts) *Service {
	return &Service{
		db:    opts.DB,
		conn:  opts.Conn,
		alias: opts.Alias,
	}
}

type CreateParams struct {
	ID     uuid.UUID
	Url    string
	Domain string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Url, error) {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("start db transaction: %w", err)
	}
	defer tx.Rollback()

	aliasSvc := s.alias.WithTx(tx)

	alias, err := aliasSvc.GetFree(ctx)
	if err != nil {
		return nil, fmt.Errorf("could reserve free alias: %w", err)
	}
	slog.Debug("retieved free alias", "alias", alias)
	slog.Debug("marking alias as used")
	if err := aliasSvc.MarkUsed(ctx, alias); err != nil {
		return nil, err
	}

	url, err := s.db.WithTx(tx).CreateUrl(ctx, queries.CreateUrlParams{
		ID:     params.ID.UUID(),
		Alias:  alias,
		Url:    params.Url,
		Domain: params.Domain,
	})
	if err != nil {
		return nil, fmt.Errorf("store url: %w", err)
	}

	slog.Debug("comitting transaction")
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit insert url: %w", err)
	}

	return mapUrl(&queries.Url{
		ID:     url.ID,
		Alias:  url.Alias,
		Url:    url.Url,
		Domain: url.Domain,
	}), nil
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

func (s *Service) Count(ctx context.Context) (int, error) {
	count, err := s.db.CountUrls(ctx)
	if err != nil {
		return 0, fmt.Errorf("count urls: %w", err)
	}
	return int(count), nil
}

var _ Urls = &Service{}
