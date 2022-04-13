package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wilhelm-murdoch/go-stash/cmd/stash/client"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/queries"
)

var (
	username = "atapas"
	hostname = "blog.greenroots.info"

	chanPostIngest = make(chan string)
	chanFinished   = make(chan bool)
)

func Messenger() {
	for {
		select {
		case slug := <-chanPostIngest:
			go ingest.Posts.Get(slug, hostname)
		case <-chanFinished:
			close(chanFinished)
			close(chanPostIngest)
			return
		}
	}
}

func main() {
	defer func() {
		chanFinished <- true
	}()

	client := client.New()

	rewind, err := time.ParseDuration("5000h")
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	since := time.Now().Add(-rewind)

	go Messenger()

	currentPage := 0
	for {
		result, err := client.Execute(queries.New("GetTimeline", queries.GetTimeline, queries.TimelineUnmarshaler, username, currentPage))
		if err != nil {
			log.Fatal(err)
		}

		if len(result.([]queries.Post)) < 6 || len(result.([]queries.Post)) == 0 {
			break
		}

		// Search publication for any posts that have been added, or updated,
		// between now and `since`. All results are dispatched to the relevant
		// ingestion handler via channel `chanPostIngest`:
		for _, post := range result.([]queries.Post) {
			dateAdded, _ := time.Parse(time.RFC3339, post.DateAdded)
			dateUpdated, _ := time.Parse(time.RFC3339, post.DateUpdated)

			if dateAdded.After(since) || dateUpdated.After(since) {
				chanPostIngest <- post.Slug
			}
		}

		currentPage++
	}

	fmt.Printf("ingest.Posts.Length(): %v\n", ingest.Posts.Length())

	// optional sleep between pages
	// save articles as json
	// json file name == ultimate url
	// adding .json to any url should display actual json file
	// can be done concurrently; go-batch?
	// download all article images to their own dedicated folder "/images/article-slug-name/<image>.png"
	//   same with images in article content

}
