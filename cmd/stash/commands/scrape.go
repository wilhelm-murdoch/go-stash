package commands

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/queries"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

// ScrapeHandler is responsible for fetching content from the specified
// Hashnode website and author:
func ScrapeHandler(c *cli.Context, cfg *config.Configuration) error {
	client := client.New()

	var since time.Time
	since = time.Now()
	if c.String("since") != "" {
		rewind, err := time.ParseDuration(c.String("since"))
		if err != nil {
			return fmt.Errorf("expected --since value to use format of 10s, 10m or 10h, but got `%s` instead", c.String("since"))
		}
		log.Printf("fetching content from `%s` ago due to usage of --since\n", c.String("since"))
		since = since.Add(-rewind)
	} else {
		rewind, _ := time.ParseDuration("99999h")
		log.Println("--since flag not used, so fetching content from the beginning")
		since = since.Add(-rewind)
	}

	currentPage := 0
	wg := new(sync.WaitGroup)
	for {
		result, err := client.Execute(queries.New("GetTimeline", queries.GetTimeline, queries.TimelineUnmarshaler, c.String("username"), currentPage))
		if err != nil {
			return err
		}

		if len(result.([]queries.Post)) == 0 {
			log.Println("Done paging.")
			break
		}
		log.Printf("Searching Page: %d\n", currentPage+1)

		// Search publication for any posts that have been added, or updated,
		// between now and `since`. All results are sent to the post ingestion
		// handler:
		for _, post := range result.([]queries.Post) {
			dateAdded, _ := time.Parse(time.RFC3339, post.DateAdded)
			dateUpdated, _ := time.Parse(time.RFC3339, post.DateUpdated)

			if dateAdded.After(since) || dateUpdated.After(since) {
				wg.Add(1)
				log.Printf("Found Article: %s", post.Slug)
				go ingest.Posts.Get(post.Slug, c.String("hostname"), wg)
			}
		}

		currentPage++
	}
	wg.Wait()

	basePathArticles := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Articles)
	if err := writers.Write(basePathArticles, ingest.Posts.Results().Items(), false); err != nil {
		return err
	}

	if err := writers.Write(basePathArticles, ingest.Posts.GetPostSummaries(), true); err != nil {
		return err
	}

	basePathTags := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags)
	if err := writers.Write(basePathTags, ingest.Posts.GroupPostsByTag(true), false); err != nil {
		return err
	}

	if err := writers.Write(basePathTags, ingest.Posts.GroupPostsByTag(true), true); err != nil {
		return err
	}

	basePathAuthors := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors)
	if err := writers.Write(basePathAuthors, ingest.Posts.FilterDistinctAuthors(), false); err != nil {
		return err
	}

	if err := writers.Write(basePathAuthors, ingest.Posts.FilterDistinctAuthors(), true); err != nil {
		return err
	}

	return nil
}
