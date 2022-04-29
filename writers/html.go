package writers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/wilhelm-murdoch/go-collection"
	"github.com/wilhelm-murdoch/go-stash/config"
	"github.com/wilhelm-murdoch/go-stash/utils"
)

func WriteHtml(basePath string, mapping *config.Mapping, data map[string]any, cfg *config.Configuration) error {
	mapping.Partials = append(mapping.Partials, mapping.Input)

	templates, err := template.New("").Funcs(sprig.FuncMap()).ParseFiles(mapping.Partials...)
	if err != nil {
		return err
	}

	f, err := os.Create(basePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data["config"] = cfg

	var buffer strings.Builder
	if err := templates.ExecuteTemplate(&buffer, string(mapping.Type), data); err != nil {
		return err
	}

	f.WriteString(buffer.String())
	log.Printf("wrote %s", basePath)

	return nil
}

func WriteHtmlCollection[W Writable](basePath string, items *collection.Collection[W], mapping *config.Mapping, data map[string]any, cfg *config.Configuration) error {
	items.Each(func(i int, item W) bool {
		if slug, err := utils.GetSlugFromItem(item); err == nil {
			data[string(mapping.Type)] = item

			if err := WriteHtml(fmt.Sprintf("%s/%s/index.html", basePath, slug), mapping, data, cfg); err != nil {
				log.Print(err)
			}
		}
		return false
	})

	return nil
}
