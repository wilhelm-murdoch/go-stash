package models

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/wilhelm-murdoch/go-stash/config"
)

type Author struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Tagline     string `json:"tagLine"`
	Photo       string `json:"photo"`
	SocialMedia struct {
		Twitter       string `json:"twitter"`
		Github        string `json:"github"`
		Stackoverflow string `json:"stackoverflow"`
		Linkedin      string `json:"linkedin"`
		Google        string `json:"google"`
		Website       string `json:"website"`
		Facebook      string `json:"facebook"`
	} `json:"socialMedia"`
}

func (a Author) GetUrl(cfg *config.Configuration) string {
	return fmt.Sprintf("%s/author/%s", cfg.Url, a.GetSlug())
}

func (a Author) GetDateUpdated() string {
	return time.Now().Format(time.RFC3339)
}

func (a Author) GetSlug() string {
	return strings.ToLower(a.Username)
}

func (a Author) GetImages(cfg *config.Configuration) []Image {
	return []Image{{a.Photo, fmt.Sprintf("%s/%s/images/author-%s-%s", cfg.Paths.Root, cfg.Paths.Files, a.GetSlug(), filepath.Base(a.Photo))}}
}

func (a *Author) ReplaceImagePaths(cfg *config.Configuration) {
	a.Photo = fmt.Sprintf("/%s/images/author-%s-%s", cfg.Paths.Files, a.GetSlug(), filepath.Base(a.Photo))
}
