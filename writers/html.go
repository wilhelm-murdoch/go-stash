package writers

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func WriteHtmlSingle(mapping *config.Mapping, data map[string]any, cfg *config.Configuration) error {
	mapping.Partials = append(mapping.Partials, mapping.Input)

	templates, err := template.New("").Funcs(sprig.FuncMap()).ParseFiles(mapping.Partials...)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", cfg.Paths.Root, mapping.Output))
	if err != nil {
		return err
	}
	defer f.Close()

	var buffer strings.Builder
	if err := templates.ExecuteTemplate(&buffer, string(mapping.Type), data); err != nil {
		return err
	}

	f.WriteString(buffer.String())

	return nil
}
