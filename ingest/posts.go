package ingest

import (
	"log"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/queries"
)

var Posts = &PostIngester{
	client:  client.New(),
	results: collection.New[*models.Post](),
}

// PostIngester
type PostIngester struct {
	client  *client.Client
	results *collection.Collection[*models.Post]
}

// Get
func (p *PostIngester) Get(slug, hostname string, wg *sync.WaitGroup) {
	defer wg.Done()

	result, err := p.client.Execute(queries.New("GetPostDetail", queries.GetPostDetail, queries.PostUnmarshaler, slug, hostname))
	if err != nil {
		log.Printf("Error Processing: %s, %s\n", slug, err.Error())
		return
	}

	log.Printf("processed post: %s\n", slug)

	p.results.Push(result.(*models.Post))
}

// Length
func (p *PostIngester) Length() int {
	return p.results.Length()
}

// Results
func (p *PostIngester) Results() *collection.Collection[*models.Post] {
	return p.results
}

// FilterDistinctAuthors
func (p *PostIngester) FilterDistinctAuthors() *collection.Collection[*models.Author] {
	authors := collection.New[*models.Author]()

	p.results.Each(func(i int, post *models.Post) bool {
		if found := authors.Contains(post.Author); !found {
			authors.Push(post.Author)
		}
		return false
	})

	return authors
}

// GroupPostsByTag
func (p *PostIngester) GroupPostsByTag() *collection.Collection[models.Tag] {
	tags := make([]models.Tag, 0)

	contains := func(needle models.Tag, haystack []models.Tag) bool {
		for _, tag := range haystack {
			if needle.Slug == tag.Slug {
				return true
			}
		}
		return false
	}

	p.results.Each(func(i int, post *models.Post) bool {
		for _, tag := range post.Tags {
			if !contains(tag, tags) {
				tags = append(tags, tag)
			}
		}
		return false
	})

	for i, tag := range tags {
		p.results.Each(func(_ int, post *models.Post) bool {
			if contains(tag, post.Tags) {
				tags[i].Posts = append(tags[i].Posts, post)
				tags[i].Count++
			}
			return false
		})
	}

	return collection.New(tags...)
}
