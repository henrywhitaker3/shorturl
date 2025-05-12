package urls

import (
	"fmt"

	"github.com/henrywhitaker3/shorturl/database/queries"
	"github.com/henrywhitaker3/shorturl/internal/uuid"
)

type Url struct {
	ID       uuid.UUID `json:"id"`
	Alias    string    `json:"alias"`
	Url      string    `json:"url"`
	ShortUrl string    `json:"short_url"`
}

func mapUrl(u *queries.Url) *Url {
	return &Url{
		ID:       uuid.UUID(u.ID),
		Alias:    u.Alias,
		Url:      u.Url,
		ShortUrl: fmt.Sprintf("https://%s/%s", u.Domain, u.Alias),
	}
}
