package models

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
