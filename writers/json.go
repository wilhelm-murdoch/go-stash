package writers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/models"
)

type Writable interface {
	models.Author | models.Post | models.Tag
}

func WriteJsonManifest[W Writable](basePath string, items *collection.Collection[W]) error {
	if err := writeJson(basePath, items.Items()); err != nil {
		return err
	}
	log.Printf("wrote %s/index.json\n", basePath)
	return nil
}

func WriteJsonBulk[W Writable](basePath string, items *collection.Collection[W]) error {
	wg := new(sync.WaitGroup)

	var slug string

	wg.Add(items.Length())
	items.Each(func(i int, item W) bool {
		go func(item W) {
			switch t := any(item).(type) {
			case models.Author:
				slug = strings.ToLower(t.Username)
			case models.Tag:
				slug = t.Slug
			case models.Post:
				slug = t.Slug
			default:
				log.Fatal("could not determine slug")
			}

			path := fmt.Sprintf("%s/%s", basePath, slug)
			defer func() {
				wg.Done()
				log.Printf("wrote %s/index.json\n", path)
			}()
			if err := writeJson(path, item); err != nil {
				log.Fatal(err)
			}
		}(item)

		return false
	})
	wg.Wait()

	return nil
}

// Save
func writeJson(path string, object any) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/index.json", path))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
