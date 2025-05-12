package urls

import (
	"fmt"

	"github.com/henrywhitaker3/go-template/database/queries"
)

type Generator interface {
	Generate() (string, error)
}

type AliasGeneratorOpts struct {
	DB *queries.Queries
}

type AliasGenerator struct {
	db *queries.Queries
}

func NewAliasGenerator(opts AliasGeneratorOpts) *AliasGenerator {
	return &AliasGenerator{
		db: opts.DB,
	}
}

func (a *AliasGenerator) Generate() (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

var _ Generator = &AliasGenerator{}
