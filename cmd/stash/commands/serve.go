package commands

import (
	"fmt"
	"net/http"

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
		http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s/index.html", cfg.Paths.Root, cfg.Paths.Posts, p.ByName("slug")))
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServePort), router)
}
