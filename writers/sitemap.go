package writers

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"strings"

	"github.com/wilhelm-murdoch/go-stash/models"
)

type Sitemap struct {
	UrlSet []SitemapUrl
}

type SitemapUrl struct {
	Loc     string
	LastMod string
}

func NewSitemap[B models.Bloggable]() *Sitemap {
	return &Sitemap{}
}

func (s *Sitemap) AddUrl(url string, lastMod string) {
	s.UrlSet = append(s.UrlSet, SitemapUrl{url, lastMod})
}

func (s *Sitemap) Save(path string) error {
	sitemap, err := template.New("").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">{{ range . }}
	<url>
		<loc>{{ .Loc }}</loc>{{ if .LastMod }}
		<lastmod>{{ .LastMod }}</lastmod>{{ end }}
	</url>{{ end }}
</urlset>`)

	if err != nil {
		return err
	}

	var buffer strings.Builder
	if sitemap.Execute(&buffer, s.UrlSet); err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/sitemap.xml", path))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(html.UnescapeString(buffer.String()))
	if err != nil {
		return err
	}

	return nil
}
