package writers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/queries"
)

func WriteFeeds(cfg *config.Configuration, items *collection.Collection[queries.Post]) error {
	feed := &feeds.Feed{
		Title:       cfg.Title,
		Link:        &feeds.Link{Href: cfg.Url},
		Description: cfg.Description,
		Created:     time.Now(),
	}

	items.Sort(func(i, j int) bool {
		l, _ := items.At(i)
		r, _ := items.At(j)

		lDate, _ := time.Parse(time.RFC3339, l.DateAdded)
		rDate, _ := time.Parse(time.RFC3339, r.DateAdded)

		return lDate.Unix() > rDate.Unix()
	}).Slice(0, cfg.FeedLimit).Each(func(i int, p queries.Post) bool {
		dateAdded, _ := time.Parse(time.RFC3339, p.DateAdded)
		dateUpdated, _ := time.Parse(time.RFC3339, p.DateUpdated)

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/%s", cfg.Url, p.Slug)},
			Description: p.Brief,
			Author:      &feeds.Author{Name: p.Author.Name},
			Created:     dateAdded,
			Updated:     dateUpdated,
		})
		return false
	})

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	feeds := map[string]string{
		"rss.xml":  rss,
		"atom.xml": atom,
	}

	for k, v := range feeds {
		file, err := os.Create(fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Feeds, k))
		if err != nil {
			return err
		}
		defer file.Close()

		file.WriteString(v)
		log.Printf("wrote %s/%s/%s", cfg.Paths.Root, cfg.Paths.Feeds, k)
	}

	return nil
}
