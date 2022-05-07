package models

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wilhelm-murdoch/go-stash/config"
)

const imageRegexPattern = `<img[^>]+\bsrc=["']([^"'?]+)["']`

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
	Author      *Author
}

func (p Post) GetUrl(cfg *config.Configuration) string {
	return fmt.Sprintf("%s/%s/%s", cfg.Url, cfg.Paths.Posts, p.GetSlug())
}

func (p Post) GetDateUpdated() string {
	if p.DateUpdated == "" {
		return p.DateAdded
	}
	return p.DateUpdated
}

func (p Post) GetSlug() string {
	return p.Slug
}

func (p Post) GetImages(cfg *config.Configuration) []Image {
	var images []Image

	destination := fmt.Sprintf("%s/%s/images", cfg.Paths.Root, cfg.Paths.Files)

	images = append(images, Image{p.CoverImage, fmt.Sprintf("%s/cover-%s-%s", destination, p.GetSlug(), filepath.Base(p.CoverImage))})

	pattern := regexp.MustCompile(imageRegexPattern)
	matches := pattern.FindAllStringSubmatch(p.Content, -1)
	for _, match := range matches {
		images = append(images, Image{match[1], fmt.Sprintf("%s/post-%s-%s", destination, p.GetSlug(), filepath.Base(match[1]))})
	}

	return images
}

func (p *Post) ReplaceImagePaths(cfg *config.Configuration) {
	p.CoverImage = fmt.Sprintf("/%s/images/cover-%s-%s", cfg.Paths.Files, p.GetSlug(), filepath.Base(p.CoverImage))
	pattern := regexp.MustCompile(imageRegexPattern)
	matches := pattern.FindAllStringSubmatch(p.Content, -1)
	for _, match := range matches {
		p.Content = strings.Replace(p.Content, match[1], fmt.Sprintf("/%s/images/post-%s-%s", cfg.Paths.Files, p.GetSlug(), filepath.Base(match[1])), -1)
	}
	p.Author.ReplaceImagePaths(cfg)
}
