package models

import "github.com/wilhelm-murdoch/go-stash/config"

type Bloggable interface {
	GetSlug() string
	GetImages(*config.Configuration) []Image
	GetDateUpdated() string
	GetUrl(*config.Configuration) string
}

type Fileable interface {
	GetOrigin() string
	GetDestination() string
}
