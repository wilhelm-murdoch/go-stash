package writers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/models"
	"golang.org/x/sync/errgroup"
)

// WriteFile
func WriteFile[F models.Fileable](file F) error {
	if _, err := os.Stat(file.GetDestination()); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(file.GetDestination()), os.ModePerm); err != nil {
		return err
	}

	f, _ := os.Create(file.GetDestination())
	defer f.Close()

	response, err := http.Get(file.GetOrigin())
	if err != nil {
		return fmt.Errorf("download failed for %s: %s", file.GetOrigin(), err)
	}
	defer response.Body.Close()

	if _, err := io.Copy(f, response.Body); err != nil {
		return fmt.Errorf("download failed for %s: %s", file.GetOrigin(), err)
	}

	return nil
}

// WriteFileCollection
func WriteFileCollection[F models.Fileable](files *collection.Collection[F]) error {
	errors := new(errgroup.Group)

	files.Each(func(i int, file F) bool {
		errors.Go(func() error {
			if err := WriteFile(file); err != nil {
				return err
			}
			log.Println("processed", file.GetDestination())
			return nil
		})

		return false
	})

	if err := errors.Wait(); err != nil {
		return err
	}

	return nil
}
