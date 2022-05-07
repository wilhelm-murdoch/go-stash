package commands

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/client"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/ingest"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/queries"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

// SyncHandler is responsible for fetching content from the specified
// Hashnode website and author:
func SyncHandler(c *cli.Context, cfg *config.Configuration) error {
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

		if len(result.([]models.Post)) == 0 {
			log.Println("done paging")
			break
		}
		log.Printf("searching page: %d\n", currentPage+1)

		// Search publication for any posts that have been added, or updated,
		// between now and `since`. All results are sent to the post ingestion
		// handler:
		for _, post := range result.([]models.Post) {
			dateAdded, _ := time.Parse(time.RFC3339, post.DateAdded)
			dateUpdated, _ := time.Parse(time.RFC3339, post.DateUpdated)

			if dateAdded.After(since) || dateUpdated.After(since) {
				wg.Add(1)
				log.Printf("found post: %s", post.Slug)
				go ingest.Posts.Get(post.Slug, c.String("hostname"), wg)
			}
		}

		currentPage++
	}
	wg.Wait()

	images := collection.New[models.Image]()

	ingest.Posts.Results().Each(func(i int, p *models.Post) bool {
		images.Push(p.GetImages(cfg)...)
		images.PushDistinct(p.Author.GetImages(cfg)...)
		return false
	})

	if err := writers.WriteFileCollection(images); err != nil {
		return err
	}

	ingest.Posts.Results().Each(func(i int, p *models.Post) bool {
		p.ReplaceImagePaths(cfg)
		return false
	})

	basePathPosts := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts)
	posts := ingest.Posts.Results()
	if err := writers.WriteJsonCollection(basePathPosts, posts); err != nil {
		return err
	}

	if err := writers.WriteJsonManifest(basePathPosts, posts); err != nil {
		return err
	}

	basePathTags := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags)
	tags := ingest.Posts.GroupPostsByTag()
	if err := writers.WriteJsonCollection(basePathTags, tags); err != nil {
		return err
	}

	if err := writers.WriteJsonManifest(basePathTags, tags); err != nil {
		return err
	}

	basePathAuthors := fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors)
	authors := ingest.Posts.FilterDistinctAuthors()
	if err := writers.WriteJsonCollection(basePathAuthors, authors); err != nil {
		return err
	}

	if err := writers.WriteJsonManifest(basePathAuthors, authors); err != nil {
		return err
	}

	return nil
}
