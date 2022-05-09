package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/template"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

// RenderHandler
func RenderHandler(c *cli.Context, cfg *config.Configuration) error {
	if err := render(c, cfg); err != nil {
		return err
	}

	if !c.Bool("watch") {
		return nil
	}

	return watch(c, cfg)
}

// UnmarshalWalkCollection
func UnmarshalWalkCollection[B models.Bloggable](basePath string) (*collection.Collection[B], error) {
	items := collection.New[B]()

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".json") && err == nil && info.Name() != "index.json" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var item B
			if err := json.Unmarshal(content, &item); err != nil {
				return err
			}

			if err != nil {
				return fmt.Errorf("failed to unmarshal `%s`: %s", path, err)
			}

			items.Push(item)
		}
		return nil
	})

	return items, err
}

// watch
func watch(c *cli.Context, cfg *config.Configuration) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	var (
		quit   = make(chan bool)
		done   = make(chan bool)
		errors = make(chan error)
	)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Chmod != fsnotify.Chmod:
					if err := render(c, cfg); err != nil {
						errors <- err
					}
				}
			case err := <-watcher.Errors:
				errors <- err
			}
		}
	}()

	err = filepath.Walk(fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Templates), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err = watcher.Add(path); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	for {
		select {
		case err := <-errors:
			close(quit)
			return err
		case <-done:
			return nil
		}
	}
}

// render
func render(c *cli.Context, cfg *config.Configuration) error {
	log.Println("unmarshaling posts from local json store")
	posts, err := UnmarshalWalkCollection[models.Post](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts))
	if err != nil {
		return fmt.Errorf("could not unmarshal posts from local store: %s", err)
	}

	posts.Sort(func(i, j int) bool {
		l, _ := posts.At(i)
		r, _ := posts.At(j)

		lDate, _ := time.Parse(time.RFC3339, l.DateAdded)
		rDate, _ := time.Parse(time.RFC3339, r.DateAdded)

		return lDate.Unix() > rDate.Unix()
	})

	log.Println("unmarshaling tags from local json store")
	tags, err := UnmarshalWalkCollection[models.Tag](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags))
	if err != nil {
		return fmt.Errorf("could not unmarshal tags from local store: %s", err)
	}

	log.Println("unmarshaling authors from local json store")
	authors, err := UnmarshalWalkCollection[models.Author](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors))
	if err != nil {
		return fmt.Errorf("could not unmarshal authors from local store: %s", err)
	}

	log.Println("rendering rss and atom feeds")
	if err := writers.WriteFeeds(cfg, posts); err != nil {
		return fmt.Errorf("failed to write feeds to %s: %s", cfg.Paths.Root, err)
	}

	data := []template.TemplateData{
		{Name: "Posts", Data: posts.Items()},
		{Name: "Tags", Data: tags.Items()},
		{Name: "Authors", Data: authors.Items()},
		{Name: "Config", Data: cfg},
	}

	sitemap := writers.NewSitemap()

	for _, mapping := range cfg.Mappings {
		log.Printf("rendering %s mappings to html", mapping.Type)
		switch mapping.Type {
		case config.Index:
			if err := writers.WriteHtml(fmt.Sprintf("%s/index.html", cfg.Paths.Root), mapping, cfg, data...); err != nil {
				return fmt.Errorf("could not render %s mapping: %s", mapping.Type, err)
			}
			sitemap.AddUrl(fmt.Sprintf("%s/index.html", cfg.Url), time.Now().Format(time.RFC3339))
		case config.Page:
			if err := writers.WriteHtml(fmt.Sprintf("%s/%s", cfg.Paths.Root, mapping.Output), mapping, cfg, data...); err != nil {
				return fmt.Errorf("could not render %s mapping: %s", mapping.Type, err)
			}
			sitemap.AddUrl(fmt.Sprintf("%s/%s", cfg.Url, mapping.Output), time.Now().Format(time.RFC3339))
		case config.Post:
			if err := writers.WriteHtmlCollection(fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts), posts, mapping, cfg, data...); err != nil {
				return fmt.Errorf("could not render %s mapping: %s", mapping.Type, err)
			}
		case config.Tag:
			if err := writers.WriteHtmlCollection(fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags), tags, mapping, cfg, data...); err != nil {
				return fmt.Errorf("could not render %s mapping: %s", mapping.Type, err)
			}
		case config.Author:
			if err := writers.WriteHtmlCollection(fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors), authors, mapping, cfg, data...); err != nil {
				return fmt.Errorf("could not render %s mapping: %s", mapping.Type, err)
			}
		default:
			log.Printf("mapping type %s is currently not supported", mapping.Type)
		}
	}

	posts.Each(func(_ int, post models.Post) bool {
		sitemap.AddUrl(post.GetUrl(cfg), post.GetDateUpdated())
		return false
	})

	tags.Each(func(_ int, tag models.Tag) bool {
		sitemap.AddUrl(tag.GetUrl(cfg), tag.GetDateUpdated())
		return false
	})

	authors.Each(func(_ int, author models.Author) bool {
		sitemap.AddUrl(author.GetUrl(cfg), author.GetDateUpdated())
		return false
	})

	sitemap.Save(cfg.Paths.Root)
	if err := sitemap.UpsertRobots(cfg); err != nil {
		return err
	}

	return nil
}
