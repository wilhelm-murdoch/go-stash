package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

type Unmarshalable interface {
	models.Author | models.Post | models.Tag
}

func UnmarshalWalkIntoCollection[U Unmarshalable](basePath string) (*collection.Collection[U], error) {
	items := collection.New[U]()

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if fmt.Sprintf("%s/index.json", basePath) != path && err == nil && info.Name() == "index.json" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var item U
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

func UnmarshalFileIntoCollection[U Unmarshalable](fullPath string) (*collection.Collection[U], error) {
	items := collection.New[U]()

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var item []U
	if err := json.Unmarshal(content, &item); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal `%s`: %s", fullPath, err)
	}

	items.Push(item...)

	return items, nil
}

func RenderHandler(c *cli.Context, cfg *config.Configuration) error {
	log.Println("unmarshaling posts from local json store")
	posts, err := UnmarshalWalkIntoCollection[models.Post](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts))
	if err != nil {
		return fmt.Errorf("could not unmarshal posts from local store: %s", err)
	}

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
		return fmt.Errorf("failed to write feeds at `%s`: %s", fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Feeds), err)
	}

	indexMapping, ok := cfg.GetIndexMapping()
	if !ok {
		return errors.New("a single mapping of type `index` must be defined")
	}

	indexData := map[string]any{
		"config":  cfg,
		"posts":   posts.Items(),
		"tags":    tags.Items(),
		"authors": authors.Items(),
	}

	log.Println("rendering index mapping to html")
	if err := writers.WriteHtmlSingle(indexMapping, indexData, cfg); err != nil {
		return fmt.Errorf("could not render index mapping: %s", err)
	}

	// articles.Each(func(i int, p queries.Post) bool {
	// 	f, err := os.Create(fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Articles, p.Slug))
	// 	if err != nil {
	// 		log.Println(err)
	// 		return true
	// 	}
	// 	defer f.Close()

	// 	data["article"] = p

	// 	var buffer strings.Builder
	// 	if err := templates.ExecuteTemplate(&buffer, "article", data); err != nil {
	// 		log.Println(err)
	// 		return true
	// 	}

	// 	f.WriteString(buffer.String())

	// 	return false
	// })

	return nil
}
