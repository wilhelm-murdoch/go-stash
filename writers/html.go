package writers

import (
	"fmt"
	"log"

	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/models"
	"github.com/wilhelm-murdoch/go-stash/template"
)

// WriteHtml
func WriteHtml(basePath string, mapping *config.Mapping, cfg *config.Configuration, data ...template.TemplateData) error {
	data = append(data, template.TemplateData{Name: "config", Data: cfg})

	template := template.NewFromMapping(mapping, data...)
	if err := template.Save(basePath); err != nil {
		return err
	}

	return nil
}

// WriteHtmlCollection
func WriteHtmlCollection[B models.Bloggable](basePath string, items *collection.Collection[B], mapping *config.Mapping, cfg *config.Configuration, data ...template.TemplateData) error {
	items.Each(func(i int, item B) bool {
		data = append(data, template.TemplateData{Name: "type", Data: mapping.Type})
		data = append(data, template.TemplateData{Name: "id", Data: item.GetSlug()})

		if err := WriteHtml(fmt.Sprintf("%s/%s/index.html", basePath, item.GetSlug()), mapping, cfg, data...); err != nil {
			log.Print(err)
		}

		return false
	})

	return nil
}
