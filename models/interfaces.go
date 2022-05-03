package models

import "github.com/wilhelm-murdoch/go-stash/config"

type Bloggable interface {
	GetSlug() string
	GetImages(*config.Configuration) []Image
}

type Fileable interface {
	GetOrigin() string
	GetDestination() string
}
