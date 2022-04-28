package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/queries"
	"github.com/wilhelm-murdoch/go-stash/writers"
)

type Unmarshalable interface {
	queries.Author | queries.Post | queries.Tag | queries.PostSummary
}

func UnmarshalIntoCollection[U Unmarshalable](basePath string) (*collection.Collection[U], error) {
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

func RenderHandler(c *cli.Context, cfg *config.Configuration) error {
	articles, err := UnmarshalIntoCollection[queries.Post](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Articles))
	if err != nil {
		return err
	}

	// tags, err := UnmarshalIntoCollection[queries.Tag](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags))
	// if err != nil {
	// 	return err
	// }

	// authors, err := UnmarshalIntoCollection[queries.Author](fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors))
	// if err != nil {
	// 	return err
	// }

	if err := writers.WriteFeeds(cfg, articles); err != nil {
		return fmt.Errorf("failed to write feeds at `%s`: %s", fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Feeds), err)
	}

	return nil
}
