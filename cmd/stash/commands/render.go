package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/template"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

func UnmarshalWalkIntoCollection[B models.Bloggable](basePath string) (*collection.Collection[B], error) {
	items := collection.New[B]()

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if fmt.Sprintf("%s/index.json", basePath) != path && err == nil && info.Name() == "index.json" {
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

func RenderHandler(c *cli.Context, cfg *config.Configuration) error {
	log.Println("unmarshaling posts from local json store")
	posts, err := UnmarshalWalkIntoCollection[models.Post](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts))
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
	tags, err := UnmarshalWalkIntoCollection[models.Tag](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags))
	if err != nil {
		return fmt.Errorf("could not unmarshal tags from local store: %s", err)
	}

	log.Println("unmarshaling authors from local json store")
	authors, err := UnmarshalWalkIntoCollection[models.Author](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors))
	if err != nil {
		return fmt.Errorf("could not unmarshal authors from local store: %s", err)
	}

	log.Println("rendering rss and atom feeds")
	if err := writers.WriteFeeds(cfg, posts); err != nil {
		return fmt.Errorf("failed to write %s: %s", fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Feeds), err)
	}

	data := []template.TemplateData{
		{Name: "Posts", Data: posts.Items()},
		{Name: "Tags", Data: tags.Items()},
		{Name: "Authors", Data: authors.Items()},
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
