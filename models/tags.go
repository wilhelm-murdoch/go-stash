package models

type Tag struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int    `json:"count,omitempty"`
	Posts []Post `json:"posts,omitempty"`
}
