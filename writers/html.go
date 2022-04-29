package writers

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func EstimateReadingTime(text string) string {
	pattern, _ := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	words := strings.Fields(pattern.ReplaceAllString(text, ""))

	return fmt.Sprintf("%.0f mins", math.Ceil(float64(len(words))/float64(200)))
}

func WriteHtmlSingle(mapping *config.Mapping, data map[string]any, cfg *config.Configuration) error {
	mapping.Partials = append(mapping.Partials, mapping.Input)

	funcMap := sprig.FuncMap()
	funcMap["readingTime"] = EstimateReadingTime

	templates, err := template.New("").Funcs(funcMap).ParseFiles(mapping.Partials...)
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
