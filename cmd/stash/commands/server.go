package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/config"
)

var (
	serverBrowsers = map[string]string{
		"darwin":  "open",
		"windows": "start",
		"linux":   "xdg-open",
	}
)

func ServerHandler(c *cli.Context, cfg *config.Configuration) error {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if bytes, err := os.ReadFile(fmt.Sprintf("%s/404.html", cfg.Paths.Root)); err == nil {
			log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
			w.Write(bytes)
		}
	})

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", cfg.Paths.Root))
	})

	router.GET("/feeds/rss", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/rss.xml", cfg.Paths.Root))
	})

	router.GET("/feeds/atom", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/atom.xml", cfg.Paths.Root))
	})

	router.GET("/robots.txt", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/robots.txt", cfg.Paths.Root))
	})

	router.GET("/sitemap.xml", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/sitemap.xml", cfg.Paths.Root))
	})

	router.GET("/posts", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Posts))
	})

	router.GET("/tags", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Tags))
	})

	router.GET("/authors", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Authors))
	})

	router.GET("/post/:slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Posts, p.ByName("slug"))
		if strings.HasSuffix(r.RequestURI, ".json") {
			path = fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Posts, p.ByName("slug"))
		}
		log.Printf("Method: %s, URI: %s\n", r.Method, path)
		http.ServeFile(w, r, path)
	})

	router.GET("/tag/:slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Tags, p.ByName("slug"))
		if strings.HasSuffix(r.RequestURI, ".json") {
			path = fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Tags, p.ByName("slug"))
		}
		log.Printf("Method: %s, URI: %s\n", r.Method, path)
		http.ServeFile(w, r, path)
	})

	router.GET("/author/:slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Authors, p.ByName("slug"))
		if strings.HasSuffix(r.RequestURI, ".json") {
			path = fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Authors, p.ByName("slug"))
		}
		log.Printf("Method: %s, URI: %s\n", r.Method, path)
		http.ServeFile(w, r, path)
	})

	router.GET("/files/*slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Methods: %s, URI: %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Files, p.ByName("slug")))
	})

	if open, ok := serverBrowsers[runtime.GOOS]; ok {
		go func() {
			cmd := exec.Command(open, fmt.Sprintf("http://localhost:%d", cfg.ServePort))
			if _, err := cmd.CombinedOutput(); err != nil {
				log.Println(err)
			}
		}()
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServePort), router)
}
