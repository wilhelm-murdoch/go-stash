package writers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/utils"
)

type Writable interface {
	models.Author | models.Post | models.Tag
}

func WriteJsonManifest[W Writable](basePath string, items *collection.Collection[W]) error {
	if err := writeJson(basePath, items.Items()); err != nil {
		return err
	}
	return nil
}

func WriteJsonBulk[W Writable](basePath string, items *collection.Collection[W]) error {
	var wg sync.WaitGroup

	write := func(item W) {
		defer wg.Done()
		if slug, err := utils.GetSlugFromItem(item); err == nil {
			if err := writeJson(fmt.Sprintf("%s/%s", basePath, slug), item); err != nil {
				log.Fatal(err)
			}
		}
	}

	wg.Add(items.Length())
	items.Each(func(i int, item W) bool {
		go write(item)
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
	defer func() {
		log.Printf("wrote %s/index.json\n", path)
		file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
