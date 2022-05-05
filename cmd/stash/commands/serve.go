package commands

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func ServeHandler(c *cli.Context, cfg *config.Configuration) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", cfg.Paths.Root))
	})

	router.GET("/:slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// check if dedicated page
		// move on to articles if not
		if strings.HasSuffix(p.ByName("slug"), ".json") {
			http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Posts, strings.TrimSuffix(p.ByName("slug"), ".json")))
			return
		}
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Posts, p.ByName("slug")))
	})

	// router.GET("/tags.json", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	http.ServeFile(w, r, fmt.Sprintf("%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Tags))
	// })

	// router.GET("/tags", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	http.ServeFile(w, r, fmt.Sprintf("%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Tags))
	// })

	// router.GET("/tag/:slug", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	if strings.HasSuffix(p.ByName("slug"), ".json") {
	// 		http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s/index.json", cfg.Paths.Root, cfg.Paths.Tags, strings.TrimSuffix(p.ByName("slug"), ".json")))
	// 		return
	// 	}
	// 	http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Tags, p.ByName("slug")))
	// })

	log.Fatal(http.ListenAndServe(":3000", router))
	return nil
}
