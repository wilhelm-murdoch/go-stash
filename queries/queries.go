package queries

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jpillora/backoff"
)

var (
	GetPostDetail = `
	post(slug:\"%s\", hostname:\"%s\") { 
    title
    slug
    cuid
		brief
		type
    coverImage
    content
    contentMarkdown
    dateAdded
    dateUpdated
		tags {
			name
			slug
			logo
		}
    author {
      username
      name
      tagline
      location
      photo
      socialMedia {
        twitter
        github
        stackoverflow
        linkedin
        google
        website
        facebook
      }
    }
	}`

	GetTimeline = `
	user(username: \"%s\") { 
		publication { 
			posts(page:%d) { 
				slug
				dateAdded
				dateUpdated
			} 
		} 
	}`
)

// Query
type Query struct {
	Name        string
	Query       string
	args        []any
	Unmarshaler func([]byte) (any, error)
	Backoff     *backoff.Backoff
}

// New
func New(name, query string, unmarshaler func([]byte) (any, error), args ...any) Query {
	backoff := &backoff.Backoff{
		Jitter: true,
	}

	return Query{
		Name:        name,
		Query:       fmt.Sprintf(`{"operationName": "%s", "query": "query %s { %s }" }`, name, name, query),
		Unmarshaler: unmarshaler,
		args:        args,
		Backoff:     backoff,
	}
}

// String
func (q *Query) String() string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(fmt.Sprintf(q.Query, q.args...), " "))
}
