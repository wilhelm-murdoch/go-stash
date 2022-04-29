package writers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/wilhelm-murdoch/go-stash/models"
)

type Writable interface {
	models.Author | models.Post | models.Tag
}

func Write[W Writable](basePath string, items []W, writeManifest bool) error {
	wg := new(sync.WaitGroup)

	if writeManifest {
		if err := writeJson(basePath, items); err != nil {
			return err
		}
		log.Printf("wrote %s/index.json\n", basePath)
		return nil
	}

	var slug string

	wg.Add(len(items))
	for _, item := range items {
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
	}

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
