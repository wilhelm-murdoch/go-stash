package commands

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/queries"
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

	// Save each post as an individual JSON file:
	wg.Add(ingest.Posts.Length())
	ingest.Posts.Results().Each(func(i int, p queries.Post) bool {
		go func() {
			path := fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Articles, p.Slug)
			defer func() {
				wg.Done()
				log.Println("Wrote:", path)
			}()

			p.ReadingTime = queries.EstimateReadingTime(p.ContentMarkdown)

			if err := ingest.Save(path, p); err != nil {
				log.Fatal(err)
			}
		}()
		return false
	})
	wg.Wait()

	// Writes a JSON file for each tag and their associated posts:
	postsByTag := ingest.Posts.GroupPostsByTag(true)
	wg.Add(len(postsByTag))
	for _, tag := range postsByTag {
		go func(tag queries.Tag) {
			path := fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Tags, tag.Slug)
			defer func() {
				wg.Done()
				log.Println("Wrote:", path)
			}()
			if err := ingest.Save(path, tag); err != nil {
				log.Fatal(err)
			}
		}(tag)
	}
	wg.Wait()

	// Writes a JSON file containing all tags with their associated posts:
	wg.Add(2)
	go func() {
		path := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags)
		defer func() {
			wg.Done()
			log.Println("Wrote:", path)
		}()
		if err := ingest.Save(path, ingest.Posts.GroupPostsByTag(true)); err != nil {
			log.Fatal(err)
		}
	}()

	// Writes a JSON file containing all posts:
	go func() {
		path := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Articles)
		defer func() {
			wg.Done()
			log.Println("Wrote:", path)
		}()
		if err := ingest.Save(path, ingest.Posts.GetPostSummaries()); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		path := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors)
		defer func() {
			wg.Done()
			log.Println("Wrote:", path)
		}()
		if err := ingest.Save(path, ingest.Posts.FilterDistinctAuthors()); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	authors := ingest.Posts.FilterDistinctAuthors()
	wg.Add(len(authors))
	for _, author := range authors {
		path := fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Authors, strings.ToLower(author.Username))
		go func(author queries.Author) {
			defer func() {
				wg.Done()
				log.Println("Wrote:", path)
			}()
			if err := ingest.Save(path, author); err != nil {
				log.Fatal(err)
			}
		}(author)
	}
	wg.Wait()

	return nil
}
