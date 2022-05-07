package writers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/models"
	"golang.org/x/sync/errgroup"
)

func WriteJsonManifest[B models.Bloggable](basePath string, items *collection.Collection[B]) error {
	if err := writeJson(fmt.Sprintf("%s/index.json", basePath), items.Items()); err != nil {
		return err
	}
	return nil
}

func WriteJsonCollection[B models.Bloggable](basePath string, items *collection.Collection[B]) error {
	errors := new(errgroup.Group)

	items.Each(func(i int, item B) bool {
		errors.Go(func() error {
			if err := writeJson(fmt.Sprintf("%s/%s.json", basePath, item.GetSlug()), &item); err != nil {
				return err
			}
			return nil
		})

		return false
	})

	if err := errors.Wait(); err != nil {
		return err
	}

	return nil
}

// Save
func writeJson(path string, object any) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("wrote %s\n", path)
		file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
