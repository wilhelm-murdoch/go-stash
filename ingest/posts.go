package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/queries"
	"github.com/wilhelm-murdoch/go-stash/tools"
)

var (
	Posts = &PostIngester{
		client:  client.New(),
		results: collection.New[queries.Post](),
	}

	contains = func(needle queries.Tag, haystack []queries.Tag) bool {
		for _, tag := range haystack {
			if needle.Slug == tag.Slug {
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
		tools.Warning(fmt.Sprintf("Error Processing: %s, %s", slug, err.Error()))
		return
	}

	tools.Info(fmt.Sprintf("Processed Article: %s", slug))
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

// FilterDistinctTags
func (p *PostIngester) FilterDistinctTags() []queries.Tag {
	tags := make([]queries.Tag, 0)

	p.results.Each(func(i int, post queries.Post) bool {
		for _, tag := range post.Tags {
			if !contains(tag, tags) {
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
			if contains(tag, post.Tags) {
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

// Save
func Save(path string, object any) error {
	file, err := os.Create(path)
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
