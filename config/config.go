package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type TemplateMapType string

const (
	Tag    TemplateMapType = "tag"
	Post   TemplateMapType = "post"
	Author TemplateMapType = "author"
	Page   TemplateMapType = "page"
	Index  TemplateMapType = "index"
)

type Configuration struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Url         string `json:"url" yaml:"url"`
	FeedLimit   int    `json:"feedLimit" yaml:"feed_limit"`
	Author      string `json:"author" yaml:"author"`
	Paths       struct {
		Root      string `json:"root" yaml:"root"`
		Feeds     string `json:"feeds" yaml:"feeds"`
		Posts     string `json:"posts" yaml:"posts"`
		Authors   string `json:"authors" yaml:"authors"`
		Tags      string `json:"tags" yaml:"tags"`
		Templates string `json:"templates" yaml:"templates"`
	} `json:"paths" yaml:"paths"`
	Mappings []*Mapping `json:"mappings" yaml:"mappings"`
}

type Mapping struct {
	Type     TemplateMapType `json:"type" yaml:"type"`
	Pattern  string          `json:"pattern" yaml:"pattern"`
	Input    string          `json:"input" yaml:"input"`
	Output   string          `json:"output" yaml:"output"`
	Partials []string        `json:"partials" yaml:"partials"`
}

// New
func New(configPath string) (*Configuration, error) {
	config := &Configuration{}

	if configPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		configPath = fmt.Sprintf("%s/.stash.yaml", cwd)
	}

	if configPath == "" {
		return nil, errors.New("you must specify a valid configuration file")
	}

	if stats, err := os.Stat(configPath); err != nil || stats.Size() == 0 {
		return nil, fmt.Errorf("specified config file `%s` is empty or does not exist", configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config.Validate()
}

// Validate
func (c *Configuration) Validate() (*Configuration, error) {
	if err := c.validatePathsBlock(); err != nil {
		return c, err
	}

	if err := c.validateMappings(); err != nil {
		return c, err
	}

	return c, nil
}

// validatePathsBlock
func (c *Configuration) validatePathsBlock() error {
	isPathValid := func(label, path string) error {
		stats, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file system error for `%s`: %s", label, err)
		}

		if !stats.IsDir() {
			return fmt.Errorf("path `%s` defined for `%s` must be a valid directory", path, label)
		}

		if stats.Mode().Perm()&(1<<(uint(7))) == 0 {
			return fmt.Errorf("path `%s` defined for `%s` must be writable", path, label)
		}

		return nil
	}

	if err := isPathValid("paths.root", c.Paths.Root); err != nil {
		return err
	}

	checks := map[string]string{
		"paths.posts":     fmt.Sprintf("%s/%s", c.Paths.Root, c.Paths.Posts),
		"paths.tags":      fmt.Sprintf("%s/%s", c.Paths.Root, c.Paths.Tags),
		"paths.authors":   fmt.Sprintf("%s/%s", c.Paths.Root, c.Paths.Authors),
		"paths.templates": fmt.Sprintf("%s/%s", c.Paths.Root, c.Paths.Templates),
		"paths.feeds":     fmt.Sprintf("%s/%s", c.Paths.Root, c.Paths.Feeds),
	}

	for label, path := range checks {
		if err := isPathValid(label, path); err != nil {
			return err
		}
	}

	return nil
}

func (c *Configuration) validateMappings() error {
	isPathValid := func(label, path string) error {
		stats, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file system error for `%s`: %s", label, err)
		}

		if stats.IsDir() {
			return fmt.Errorf("path `%s` defined for `%s` must be a valid file", path, label)
		}

		return nil
	}

	templatePath := func(fileName string) string {
		return fmt.Sprintf("%s/%s/%s", c.Paths.Root, c.Paths.Templates, fileName)
	}

	var indexMappingDefined bool
	for i1, mapping := range c.Mappings {
		if mapping.Type == Index {
			if indexMappingDefined {
				return errors.New("only a single mapping of type `index` is allowed")
			} else {
				indexMappingDefined = true
			}
		}

		if err := isPathValid(fmt.Sprintf("mappings[%d]input", i1), templatePath(mapping.Input)); err != nil {
			return err
		}

		mapping.Input = templatePath(mapping.Input)

		for i2, partial := range mapping.Partials {
			if err := isPathValid(fmt.Sprintf("mappings[%d]partials[%d]", i1, i2), templatePath(partial)); err != nil {
				return err
			}

			mapping.Partials[i2] = templatePath(partial)
		}
	}

	if !indexMappingDefined {
		return errors.New("a single mapping of type `index` must be defined")
	}

	return nil
}

func (c *Configuration) GetFirstMappingByType(mappingType TemplateMapType) (mapping *Mapping, ok bool) {
	mappings, ok := c.GetMappingsByType(mappingType)
	if !ok {
		return mapping, ok
	}

	return mappings[0], true
}

func (c *Configuration) GetMappingsByType(mappingType TemplateMapType) (mappings []*Mapping, ok bool) {
	for _, mapping := range c.Mappings {
		if mapping.Type == mappingType {
			mappings = append(mappings, mapping)
		}
	}

	if len(mappings) == 0 {
		return mappings, ok
	}

	return mappings, true
}

func WrapWithConfig(c *cli.Context, action func(*cli.Context, *Configuration) error) error {
	config, err := New(c.String("config"))
	if err != nil {
		return err
	}

	return action(c, config)
}
