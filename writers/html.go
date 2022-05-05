package writers

import (
	"fmt"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/template"
	"golang.org/x/sync/errgroup"
)

// WriteHtml
func WriteHtml(basePath string, mapping *config.Mapping, cfg *config.Configuration, data ...template.TemplateData) error {
	data = append(data, template.TemplateData{Name: "Config", Data: cfg})

	template := template.NewFromMapping(mapping, data...)
	if err := template.Save(basePath); err != nil {
		return err
	}

	return nil
}

// WriteHtmlCollection
func WriteHtmlCollection[B models.Bloggable](basePath string, items *collection.Collection[B], mapping *config.Mapping, cfg *config.Configuration, data ...template.TemplateData) error {
	errors := new(errgroup.Group)

	items.Each(func(i int, item B) bool {
		slug := item.GetSlug()

		errors.Go(func() error {
			data = append(data, template.TemplateData{Name: "Type", Data: mapping.Type})
			data = append(data, template.TemplateData{Name: "Id", Data: slug})

			if err := WriteHtml(fmt.Sprintf("%s/%s/index.html", basePath, slug), mapping, cfg, data...); err != nil {
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
