package utils

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/wilhelm-murdoch/go-stash/models"
)

type Writable interface {
	models.Author | models.Post | models.Tag
}

func GetSlugFromItem[W Writable](item W) (string, error) {
	var slug string
	switch t := any(item).(type) {
	case models.Author:
		slug = strings.ToLower(t.Username)
	case models.Tag:
		slug = t.Slug
	case models.Post:
		slug = t.Slug
	default:
		return "", errors.New("could not determine slug")
	}

	return slug, nil
}

func EstimateReadingTime(text string) string {
	pattern, _ := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	words := strings.Fields(pattern.ReplaceAllString(text, ""))

	return fmt.Sprintf("%.0f mins", math.Ceil(float64(len(words))/float64(238)))
}
