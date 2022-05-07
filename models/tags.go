package models

import (
	"fmt"
	"time"

	"github.com/wilhelm-murdoch/go-stash/config"
)

type Tag struct {
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Count int     `json:"count,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

func (t Tag) GetUrl(cfg *config.Configuration) string {
	return fmt.Sprintf("%s/tag/%s", cfg.Url, t.GetSlug())
}

func (t Tag) GetDateUpdated() string {
	return time.Now().Format(time.RFC3339)
}

func (t Tag) GetSlug() string {
	return t.Slug
}

func (t Tag) GetImages(cfg *config.Configuration) []Image {
	return nil
}
