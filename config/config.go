package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type TemplateMapType string

const (
	Unknown TemplateMapType = "unknown"
	Tag     TemplateMapType = "tag"
	Article TemplateMapType = "article"
	Page    TemplateMapType = "page"
	Index   TemplateMapType = "index"
)

type Config struct {
	Directories struct {
		Articles  string `yaml:"articles"`
		Tags      string `yaml:"tags"`
		Authors   string `yaml:"authors"`
		Templates string `yaml:"templates"`
		Images    string `yaml:"images"`
		Static    string `yaml:"static"`
	} `yaml:"directories"`
	Index struct {
		Input    string   `yaml:"input"`
		Output   string   `yaml:"output"`
		Sources  []string `yaml:"sources"`
		Partials []string `yaml:"partials"`
	} `yaml:"index"`
}

func New(configPath string) (*Config, error) {
	config := &Config{}

	if stats, err := os.Stat(configPath); err != nil || stats.Size() == 0 {
		return nil, fmt.Errorf("specified config file `%s` is empty", configPath)
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

func (c *Config) Validate() (*Config, error) {
	if err := c.validateDirectoriesBlock(); err != nil {
		return c, err
	}

	if err := c.validateIndexBlock(); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) validateDirectoriesBlock() error {
	isValid := func(label, path string) error {
		stats, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file system error for `%s`: %s", label, err)
		}

		if !stats.IsDir() {
			return fmt.Errorf("path `%s` defined for `%s` must be a valid directory", path, label)
		}

		return nil
	}

	checks := map[string]string{
		"directories.articles":  c.Directories.Articles,
		"directories.tags":      c.Directories.Tags,
		"directories.authors":   c.Directories.Authors,
		"directories.templates": c.Directories.Templates,
		"directories.images":    c.Directories.Images,
		"directories.static":    c.Directories.Static,
	}

	for k, v := range checks {
		if err := isValid(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) validateIndexBlock() error {
	isValid := func(label, path string) error {
		stats, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file system error for `%s`: %s", label, err)
		}

		if stats.IsDir() {
			return fmt.Errorf("path `%s` defined for `%s` must be a valid file", path, label)
		}

		return nil
	}

	if err := isValid("index.input", fmt.Sprintf("%s/%s", c.Directories.Templates, c.Index.Input)); err != nil {
		return err
	}

	checks := map[string][]string{
		"index.sources": {c.Directories.Templates, c.Index.Sources},
	}

	return nil
}
