package ingest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/queries"
)

var (
	Posts = &PostIngester{
		client:  client.New(),
		results: collection.New[queries.Post](),
	}

	tagsContain = func(needle queries.Tag, haystack []queries.Tag) bool {
		for _, tag := range haystack {
			if needle.Slug == tag.Slug {
				return true
			}
		}
		return false
	}

	authorsContain = func(needle queries.Author, haystack []queries.Author) bool {
		for _, tag := range haystack {
			if needle.Username == tag.Username {
				return true
			}
		}
		return false
	}
)

// PostIngester
type PostIngester struct {
	client  *client.Client
	results *collection.Collection[queries.Post]
}

// Get
func (p *PostIngester) Get(slug, hostname string, wg *sync.WaitGroup) {
	defer wg.Done()

	result, err := p.client.Execute(queries.New("GetPostDetail", queries.GetPostDetail, queries.PostUnmarshaler, slug, hostname))
	if err != nil {
		log.Printf("Error Processing: %s, %s\n", slug, err.Error())
		return
	}

	log.Printf("Processed Article: %s\n", slug)

	p.results.Push(result.(queries.Post))
}

// Empty
func (p *PostIngester) Empty() {
	p.results.Empty()
}

// Length
func (p *PostIngester) Length() int {
	return p.results.Length()
}

// Results
func (p *PostIngester) Results() *collection.Collection[queries.Post] {
	return p.results
}

// GetPostSummaries
func (p *PostIngester) GetPostSummaries() []queries.PostSummary {
	posts := make([]queries.PostSummary, 0)

	p.results.Each(func(i int, post queries.Post) bool {
		posts = append(posts, queries.NewPostSummary(post))
		return false
	})

	return posts
}

// FilterDistinctAuthors
func (p *PostIngester) FilterDistinctAuthors() []queries.Author {
	authors := make([]queries.Author, 0)

	p.results.Each(func(i int, post queries.Post) bool {

		if !authorsContain(post.Author, authors) {
			authors = append(authors, post.Author)
		}

		return false
	})

	return authors
}

// FilterDistinctTags
func (p *PostIngester) FilterDistinctTags() []queries.Tag {
	tags := make([]queries.Tag, 0)

	p.results.Each(func(i int, post queries.Post) bool {
		for _, tag := range post.Tags {
			if !tagsContain(tag, tags) {
				tags = append(tags, tag)
			}
		}

		return false
	})

	return tags
}

// GroupPostsByTag
func (p *PostIngester) GroupPostsByTag(includePostSummary bool) []queries.Tag {
	tags := p.FilterDistinctTags()

	for i, tag := range tags {
		p.results.Each(func(_ int, post queries.Post) bool {
			if tagsContain(tag, post.Tags) {
				if includePostSummary {
					tags[i].Posts = append(tags[i].Posts, queries.NewPostSummary(post))
				}

				tags[i].Count++
			}
			return false
		})
	}

	return tags
}

func SaveEach(path string, object any) error {

	switch object.(type) {
	case []queries.Post:
		fmt.Println(reflect.TypeOf(object))
	}
	return nil
}

// Save
func Save(path string, object any) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/index.json", path))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
