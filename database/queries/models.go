// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package queries

import (
	"github.com/google/uuid"
)

type Alias struct {
	Alias string
	Used  bool
}

type Click struct {
	ID        uuid.UUID
	UrlID     uuid.UUID
	Ip        string
	ClickedAt int64
}

type Url struct {
	ID     uuid.UUID
	Alias  string
	Url    string
	Domain string
}
