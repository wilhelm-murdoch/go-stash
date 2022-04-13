package ingest

import (
	"log"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/client"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/queries"
)

var Posts = &PostIngester{
	client:  client.New(),
	results: collection.New[queries.Post](),
}

type PostIngester struct {
	client  *client.Client
	results *collection.Collection[queries.Post]
}

func (p *PostIngester) Get(slug, hostname string, wg *sync.WaitGroup) {
	defer wg.Done()

	result, err := p.client.Execute(queries.New("GetPostDetail", queries.GetPostDetail, queries.PostUnmarshaler, slug, hostname))
	if err != nil {
		log.Fatal(err)
	}

	p.results.Push(result.(queries.Post))
	// log.Println("completed ingesting:", slug)
}

func (p *PostIngester) Empty() {
	p.results.Empty()
}

func (p *PostIngester) Length() int {
	return p.results.Length()
}

func (p *PostIngester) FilterDistinctTags() *collection.Collection[*queries.Tag] {
	tags := collection.New[*queries.Tag]()

	p.results.Each(func(i int, p queries.Post) bool {
		for _, tag := range p.Tags {
			tags.PushDistinct(&tag)
		}
		return false
	})

	return tags
}

func (p *PostIngester) FilterPostsByTag() *collection.Collection[*queries.Tag] {
	tags := p.FilterDistinctTags()

	tags.Each(func(i int, t *queries.Tag) bool {
		p.results.Each(func(i int, post queries.Post) bool {
			for _, tag := range post.Tags {
				if tag.Slug == t.Slug {
					t.Posts = append(t.Posts, queries.NewPostSummary(post))
				}
			}
			return false
		})

		return false
	})

	return tags
}
