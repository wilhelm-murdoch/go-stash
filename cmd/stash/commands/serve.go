package commands

import (
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func ServeHandler(c *cli.Context, cfg *config.Configuration) error {
	http.Handle("/", http.FileServer(http.Dir("dist")))
	http.ListenAndServe(":3000", nil)
	return nil
}
