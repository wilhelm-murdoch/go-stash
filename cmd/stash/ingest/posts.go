package ingest

import (
	"fmt"
	"log"
	"time"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/client"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/queries"
)

var Posts = &PostIngester{
	client:     client.New(),
	collection: collection.New[queries.Post](),
}

type PostIngester struct {
	client     *client.Client
	collection *collection.Collection[queries.Post]
}

func (p *PostIngester) Get(slug, hostname string) {
	result, err := p.client.Execute(queries.New("GetPostDetail", queries.GetPostDetail, queries.PostUnmarshaler, slug, hostname))
	if err != nil {
		log.Fatal(err)
	}

	p.collection.Push(result.(queries.Post))

	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println(p.collection.Length(), result.(queries.Post).Title)
}

func (p *PostIngester) Empty() {
	p.collection.Empty()
}

func (p *PostIngester) Length() int {
	return p.collection.Length()
}
