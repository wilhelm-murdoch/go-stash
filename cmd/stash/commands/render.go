package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/config"
)

func RenderHandler(c *cli.Context, cfg *config.Configuration) error {
	// template for layout
	// template for tags
	// template for index
	// template for article
	fmt.Println("hi")
	return nil
}
