package template

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/wilhelm-murdoch/go-stash/config"
)

type Template struct {
	Name     string
	Input    string
	Output   string
	Partials []string
	Data     []TemplateData
}

type TemplateData struct {
	Name string
	Data any
}

func NewFromMapping(mapping *config.Mapping, data ...TemplateData) *Template {
	return &Template{
		Name:     strings.TrimSuffix(filepath.Base(mapping.Input), filepath.Ext(mapping.Input)),
		Input:    mapping.Input,
		Output:   mapping.Output,
		Partials: append(mapping.Partials, mapping.Input),
		Data:     data,
	}
}

func (t *Template) mapData() map[string]any {
	data := make(map[string]any, 0)

	for _, td := range t.Data {
		data[td.Name] = td.Data
	}

	return data
}

func (t *Template) Save(basePath string) error {
	if err := os.MkdirAll(filepath.Dir(basePath), os.ModePerm); err != nil {
		return err
	}

	funcMap := sprig.FuncMap()

	for name, fn := range funcMapStash {
		funcMap[name] = fn
	}

	templates, err := template.New("").Funcs(funcMap).ParseFiles(t.Partials...)
	if err != nil {
		return err
	}

	file, err := os.Create(basePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buffer strings.Builder
	if err := templates.ExecuteTemplate(&buffer, t.Name, t.mapData()); err != nil {
		return err
	}

	_, err = file.WriteString(buffer.String())
	if err != nil {
		return err
	}

	log.Printf("wrote %s", basePath)

	return nil
}
