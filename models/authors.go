package models

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wilhelm-murdoch/go-stash/config"
)

type Author struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Tagline     string `json:"tagLine"`
	Photo       string `json:"photo"`
	CoverImage  string `json:"coverImage"`
	SocialMedia struct {
		Twitter       string `json:"twitter"`
		Github        string `json:"github"`
		StackOverflow string `json:"stackoverflow"`
		LinkedIn      string `json:"linkedin"`
		Google        string `json:"google"`
		Website       string `json:"website"`
		Facebook      string `json:"facebook"`
	} `json:"socialMedia"`
}

func (a Author) GetUrl(cfg *config.Configuration) string {
	return fmt.Sprintf("%s/author/%s", cfg.Url, a.GetSlug())
}

func (a Author) GetDateUpdated() string {
	return ""
}

func (a Author) GetSlug() string {
	return strings.ToLower(a.Username)
}

func (a Author) GetImages(cfg *config.Configuration) []Image {
	var images []Image

	destination := fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Authors, a.GetSlug())

	if a.CoverImage != "" {
		images = append(images, Image{a.CoverImage, fmt.Sprintf("%s/%s", destination, filepath.Base(a.CoverImage))})
	}

	if a.Photo != "" {
		images = append(images, Image{a.Photo, fmt.Sprintf("%s/%s", destination, filepath.Base(a.Photo))})
	}

	return images
}
