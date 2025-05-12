package urls

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/henrywhitaker3/shorturl/database/queries"
)

type Alias struct {
	db *queries.Queries
}

type AliasOpts struct {
	DB *queries.Queries
}

func NewAlias(opts AliasOpts) *Alias {
	return &Alias{
		db: opts.DB,
	}
}

func (a *Alias) Count(ctx context.Context) (int, error) {
	count, err := a.db.CountAliases(ctx)
	if err != nil {
		return 0, fmt.Errorf("count all aliases: %w", err)
	}
	return int(count), nil
}

func (a *Alias) CountFree(ctx context.Context) (int, error) {
	count, err := a.db.CountFreeAliases(ctx)
	if err != nil {
		return 0, fmt.Errorf("count free aliases: %w", err)
	}
	return int(count), nil
}

func (a *Alias) Store(ctx context.Context, alias string) error {
	_, err := a.db.CreateAlias(ctx, alias)
	if err != nil {
		return fmt.Errorf("store alias: %w", err)
	}
	return nil
}

func (a *Alias) GetFree(ctx context.Context) (string, error) {
	alias, err := a.db.GetFreeAlias(ctx)
	if err != nil {
		return "", fmt.Errorf("return free alias: %w", err)
	}

	return alias.Alias, nil
}

func (a *Alias) MarkUsed(ctx context.Context, alias string) error {
	err := a.db.MarkAliasUsed(ctx, alias)
	if err != nil {
		return fmt.Errorf("mark alias as used: %w", err)
	}
	return nil
}

func (a *Alias) Create(ctx context.Context, alias string) error {
	_, err := a.db.CreateAlias(ctx, alias)
	if err != nil {
		return fmt.Errorf("store alias: %w", err)
	}
	return nil
}

func (a *Alias) FilterOutExisting(ctx context.Context, aliases []string) ([]string, error) {
	exists, err := a.db.GetAliases(ctx, aliases)
	if err != nil {
		return nil, fmt.Errorf("filter used aliases: %w", err)
	}

	out := []string{}
	for _, e := range aliases {
		if slices.Contains(exists, e) {
			continue
		}
		out = append(out, e)
	}

	return out, nil
}

func (a *Alias) WithTx(tx *sql.Tx) *Alias {
	return &Alias{
		db: a.db.WithTx(tx),
	}
}
