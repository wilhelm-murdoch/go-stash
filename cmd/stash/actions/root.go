package actions

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/queries"
)

// RootHandler is responsible for running the root command for the cli.
func RootHandler(c *cli.Context) error {
	client := client.New()

	rewind, err := time.ParseDuration(c.String("since"))
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	since := time.Now().Add(-rewind)

	currentPage := 0
	slugs := make([]string, 0)
	for {
		result, err := client.Execute(queries.New("GetTimeline", queries.GetTimeline, queries.TimelineUnmarshaler, c.String("username"), currentPage))
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
				slugs = append(slugs, post.Slug)
			}
		}

		currentPage++
	}

	wg := new(sync.WaitGroup)
	for _, s := range slugs {
		wg.Add(1)
		go ingest.Posts.Get(s, c.String("hostname"), wg)
	}
	wg.Wait()

	// encoder := json.NewEncoder(os.Stdout)
	// if err := encoder.Encode(ingest.Posts.FilterDistinctTags()); err != nil {
	// 	log.Fatal(err)
	// }

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(ingest.Posts.GroupPostsByTag(true)); err != nil {
		log.Fatal(err)
	}

	// optional sleep between pages
	// save articles as json
	// json file name == ultimate url
	// adding .json to any url should display actual json file
	// can be done concurrently; go-batch?
	// download all article images to their own dedicated folder "/images/article-slug-name/<image>.png"
	//   same with images in article content
	return nil
}
