package commands

import (
	"net/http"

	"github.com/urfave/cli/v2"
)

func ServeHandler(c *cli.Context) error {
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.ListenAndServe(":3000", nil)
	return nil
}
