package main

import (
	"fmt"
	"log"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/client"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/queries"
)

var (
	username = "blitzkriegpunk"
	hostname = "wilhelm.codes"
)

func main() {
	client := client.New()

	posts := collection.New[queries.Post]()
	currentPage := 0
	var total int
	for {
		result, err := client.Execute(queries.New("GetPosts", queries.GetTimeline, queries.PostsUnmarshaler, username, currentPage))
		if err != nil {
			log.Fatal(err)
		}

		if len(result.([]queries.Post)) == 0 {
			fmt.Println("... no more posts to find")
			break
		}

		total = posts.Push(result.([]queries.Post)...)

		fmt.Printf("on page %d and have a total of %d posts...\n", currentPage, total)

		currentPage++
	}

	posts.Each(func(i int, p queries.Post) bool {
		fmt.Println(p.Slug)
		return false
	})

	// optional sleep between pages
	// save articles as json
	// json file name == ultimate url
	// adding .json to any url should display actual json file
	// can be done concurrently; go-batch?
	// download all article images to their own dedicated folder "/images/article-slug-name/<image>.png"
	//   same with images in article content

}
