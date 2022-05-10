package template

import (
	"fmt"
	"html/template"
	"math"
	"regexp"
	"strings"

	"github.com/wilhelm-murdoch/go-stash/models"
)

var funcMapStash = map[string]any{
	"EstimateReadingTime": EstimateReadingTime,
	"Unescape":            Unescape,
	"GetPostById":         GetPostById,
	"GetTagById":          GetTagById,
	"GetAuthorById":       GetAuthorById,
	"GetNextPost":         GetNextPost,
	"GetPreviousPost":     GetPreviousPost,
}

func EstimateReadingTime(text string) string {
	pattern, _ := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	words := strings.Fields(pattern.ReplaceAllString(text, ""))

	return fmt.Sprintf("%.0f", math.Ceil(float64(len(words))/float64(238)))
}

func Unescape(text string) template.HTML {
	return template.HTML(text)
}

func GetPostById(posts []models.Post, id string) (models.Post, error) {
	var found models.Post
	for _, post := range posts {
		if post.GetSlug() == id {
			return post, nil
		}
	}

	return found, fmt.Errorf("post not found for %s", id)
}

func GetTagById(tags []models.Tag, id string) (models.Tag, error) {
	var found models.Tag
	for _, tag := range tags {
		if tag.GetSlug() == id {
			return tag, nil
		}
	}

	return found, fmt.Errorf("tag not found for %s", id)
}

func GetAuthorById(authors []models.Author, id string) (models.Author, error) {
	var found models.Author
	for _, author := range authors {
		if author.GetSlug() == id {
			return author, nil
		}
	}

	return found, fmt.Errorf("author not found for %s", id)
}

func GetNextPost() models.Post {
	var found models.Post
	return found
}

func GetPreviousPost() models.Post {
	var found models.Post
	return found
}
