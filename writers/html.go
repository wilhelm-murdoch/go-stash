package writers

import (
	"fmt"
	"log"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
)

func WriteHtml(basePath string, mapping *config.Mapping, cfg *config.Configuration, data ...models.TemplateData) error {
	data = append(data, models.TemplateData{Name: "config", Data: cfg})

	template := models.NewTemplateFromMapping(mapping, data...)
	if err := template.Save(basePath); err != nil {
		return err
	}

	return nil
}

func WriteHtmlCollection[B models.Bloggable](basePath string, items *collection.Collection[B], mapping *config.Mapping, cfg *config.Configuration, data ...models.TemplateData) error {
	items.Each(func(i int, item B) bool {
		data = append(data, models.TemplateData{Name: string(mapping.Type), Data: item})

		if err := WriteHtml(fmt.Sprintf("%s/%s/index.html", basePath, item.GetSlug()), mapping, cfg, data...); err != nil {
			log.Print(err)
		}

		return false
	})

	return nil
}
