package commands

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/queries"
)

// ScrapeHandler is responsible for running the root command for the cli.
func ScrapeHandler(c *cli.Context) error {
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

	wg.Add(ingest.Posts.Length())
	ingest.Posts.Results().Each(func(i int, p queries.Post) bool {
		go func() {
			defer func() {
				wg.Done()
				log.Printf("Wrote: dist/%s.json\n", p.Slug)
			}()

			p.ReadingTime = queries.EstimateReadingTime(p.ContentMarkdown)

			if err := ingest.Save(fmt.Sprintf("dist/%s.json", p.Slug), p); err != nil {
				log.Fatal(err)
			}
		}()
		return false
	})
	wg.Wait()

	postsByTag := ingest.Posts.GroupPostsByTag(true)
	wg.Add(len(postsByTag))
	for _, tag := range postsByTag {
		go func(tag queries.Tag) {
			defer func() {
				wg.Done()
				log.Printf("Wrote: dist/tags/%s.json\n", tag.Slug)
			}()
			if err := ingest.Save(fmt.Sprintf("dist/tags/%s.json", tag.Slug), tag); err != nil {
				log.Fatal(err)
			}
		}(tag)
	}
	wg.Wait()

	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
			log.Println("Wrote: dist/tags/index.json")
		}()
		if err := ingest.Save("dist/tags/index.json", ingest.Posts.GroupPostsByTag(true)); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			log.Println("Wrote: dist/index.json")
		}()
		if err := ingest.Save("dist/index.json", ingest.Posts.GetPostSummaries()); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	return nil
}
