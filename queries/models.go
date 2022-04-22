package queries

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// Author
type Author struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Tagline     string `json:"tagLine"`
	Photo       string `json:"photo"`
	CoverImage  string `json:"coverImage"`
	SocialMedia SocialMedia
}

// SocialMedia
type SocialMedia struct {
	Twitter       string `json:"twitter"`
	Github        string `json:"github"`
	StackOverflow string `json:"stackoverflow"`
	LinkedIn      string `json:"linkedin"`
	Google        string `json:"google"`
	Website       string `json:"website"`
	Facebook      string `json:"facebook"`
}

// Tag
type Tag struct {
	Name  string        `json:"name"`
	Slug  string        `json:"slug"`
	Count int           `json:"count,omitempty"`
	Posts []PostSummary `json:"posts,omitempty"`
}

// Post
type Post struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	DateAdded       string `json:"dateAdded"`
	DateUpdated     string `json:"dateUpdated"`
	CUID            string `json:"cuid"`
	Type            string `json:"type"`
	Content         string `json:"content"`
	ContentMarkdown string
	Brief           string `json:"brief"`
	CoverImage      string `json:"coverImage"`
	ReadingTime     string `json:"readingTime"`
	Tags            []Tag  `json:"tags,omitempty"`
	Author          Author
}

type PostSummary struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Brief       string `json:"brief"`
	CoverImage  string `json:"coverImage"`
	ReadingTime string `json:"readingTime"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Photo       string `json:"photo"`
	Tags        []Tag  `json:"tags,omitempty"`
}

func NewPostSummary(p Post) PostSummary {
	return PostSummary{
		Title:       p.Title,
		Slug:        p.Slug,
		Brief:       p.Brief,
		CoverImage:  p.CoverImage,
		ReadingTime: EstimateReadingTime(p.ContentMarkdown),
		Username:    p.Author.Username,
		Name:        p.Author.Name,
		Photo:       p.Author.Photo,
		Tags:        p.Tags,
	}
}

func EstimateReadingTime(text string) string {
	pattern, _ := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	words := strings.Fields(pattern.ReplaceAllString(text, ""))

	return fmt.Sprintf("%.0f mins", math.Ceil(float64(len(words))/float64(200)))
}
