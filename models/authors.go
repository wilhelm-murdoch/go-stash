package models

import "strings"

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

func (a Author) GetSlug() string {
	return strings.ToLower(a.Username)
}
