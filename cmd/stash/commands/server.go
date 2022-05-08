package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

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
	var (
		serverRoutes = map[string]string{
			"/":             fmt.Sprintf("%s/index.html", cfg.Paths.Root),
			"/feeds/rss":    fmt.Sprintf("%s/rss.xml", cfg.Paths.Root),
			"/feeds/atom":   fmt.Sprintf("%s/atom.xml", cfg.Paths.Root),
			"/robots.txt":   fmt.Sprintf("%s/robots.txt", cfg.Paths.Root),
			"/sitemap.xml":  fmt.Sprintf("%s/sitemap.xml", cfg.Paths.Root),
			"/authors.json": fmt.Sprintf("%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Authors),
			"/tags.json":    fmt.Sprintf("%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Tags),
			"/posts.json":   fmt.Sprintf("%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Posts),
		}

		serverRoutesSlug = map[string]string{
			"/author/:slug": fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Authors),
			"/tag/:slug":    fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Tags),
			"/post/:slug":   fmt.Sprintf("%s/%s", cfg.Paths.Root, cfg.Paths.Posts),
		}
	)

	log.Printf("starting local server on port :%d", cfg.ServePort)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if bytes, err := os.ReadFile(fmt.Sprintf("%s/404.html", cfg.Paths.Root)); err == nil {
			log.Printf("%s %s\n", r.Method, r.RequestURI)
			w.Write(bytes)
		}
	})

	for route, file := range serverRoutes {
		r, f := route, file // NFI why, but these 2 values get randomised with each iteration unless I do this
		log.Printf("... registering route %s for %s", r, f)
		router.GET(r, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			log.Printf("%s %s\n", r.Method, r.RequestURI)
			http.ServeFile(w, r, f)
		})
	}

	for route, file := range serverRoutesSlug {
		r, f := route, file // NFI why, but these 2 values get randomised with each iteration unless I do this
		log.Printf("... registering route %s for %s", r, f)
		router.GET(r, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			path := fmt.Sprintf("%s/%s/index.html", f, p.ByName("slug"))
			if strings.HasSuffix(r.RequestURI, ".json") {
				path = fmt.Sprintf("%s/%s", f, p.ByName("slug"))
			}
			log.Printf("%s %s (%s)\n", r.Method, r.RequestURI, path)
			http.ServeFile(w, r, path)
		})
	}

	for _, mapping := range cfg.Mappings {
		if mapping.Type == config.Page {
			route := fmt.Sprintf("/%s", strings.TrimSuffix(filepath.Base(mapping.Input), filepath.Ext(mapping.Input)))
			file := fmt.Sprintf("%s/%s", cfg.Paths.Root, mapping.Output)
			log.Printf("... registering mapping route %s for %s", route, file)
			router.GET(route, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				log.Printf("%s %s\n", r.Method, r.RequestURI)
				http.ServeFile(w, r, file)
			})
		}
	}

	log.Printf("... registering route /files/*slug for %s/%s/*slug", cfg.Paths.Root, cfg.Paths.Files)
	router.GET("/files/*slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s", cfg.Paths.Root, cfg.Paths.Files, p.ByName("slug")))
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServePort), Handler: router}

	go server.ListenAndServe()

	if open, ok := serverBrowsers[runtime.GOOS]; ok && !c.Bool("no-browser") {
		go func() {
			cmd := exec.Command(open, fmt.Sprintf("http://localhost:%d/%s", cfg.ServePort, c.String("path")))
			if _, err := cmd.CombinedOutput(); err != nil {
				log.Println(err)
			}
		}()
	} else {
		log.Println("... skipped opening browser")
	}

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	fmt.Println()
	log.Println("shutting server down")

	return nil
}
