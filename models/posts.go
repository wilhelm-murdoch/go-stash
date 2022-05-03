package models

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/wilhelm-murdoch/go-stash/config"
)

type Post struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	DateAdded   string `json:"dateAdded"`
	DateUpdated string `json:"dateUpdated"`
	CUID        string `json:"cuid"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Brief       string `json:"brief"`
	CoverImage  string `json:"coverImage"`
	Tags        []Tag  `json:"tags,omitempty"`
	Author      Author
}

func (p Post) GetSlug() string {
	return p.Slug
}

func (p Post) GetImages(cfg *config.Configuration) []Image {
	var images []Image

	destination := fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Posts, p.GetSlug())

	images = append(images, Image{p.CoverImage, fmt.Sprintf("%s/%s", destination, filepath.Base(p.CoverImage))})
	pattern := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"'?]+)["']`)

	matches := pattern.FindAllStringSubmatch(p.Content, -1)
	for _, match := range matches {
		images = append(images, Image{match[1], fmt.Sprintf("%s/%s", destination, filepath.Base(match[1]))})
	}

	return images
}
