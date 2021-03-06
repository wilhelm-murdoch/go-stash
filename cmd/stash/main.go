package main

// implement proper logger
// document
// test

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/wilhelm-murdoch/go-stash/cmd/stash/commands"
	"github.com/wilhelm-murdoch/go-stash/config"
)

var (
	// Version describes the version of the current build.
	Version = "dev"

	// Commit describes the commit of the current build.
	Commit = "none"

	// Date describes the date of the current build.
	Date = "unknown"

	// Release describes the stage of the current build, eg; development, production, etc...
	Stage = "unknown"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version: %s, Stage: %s, Commit: %s, Date: %s\n", Version, Stage, Commit, Date)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "print only the version",
		Aliases: []string{"v"},
	}

	app := &cli.App{
		Name:     "stash",
		Usage:    "a static site generator for Hashnode content",
		Version:  Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{{
			Name:  "Wilhelm Murdoch",
			Email: "wilhelm@devilmayco.de",
		}},
		Copyright: "(c) 2022 Wilhelm Codes ( https://wilhelm.codes )",
		Commands: []*cli.Command{
			{
				Name:  "sync",
				Usage: "fetches remote content from Hashnode's API and saves it locally",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Usage:   "path to the current project's `.stash.yml` configuration",
						EnvVars: []string{"CONFIG"},
					},
					&cli.StringFlag{
						Name:     "username",
						Usage:    "the @username of the target Hashnode user.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "hostname",
						Usage:    "the hostname of the target Hashnode blog.",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "since",
						Usage: "return content that occured since this period, eg; 10m, 10h",
					},
				},
				Action: func(c *cli.Context) error {
					return config.WrapWithConfig(c, commands.SyncHandler)
				},
			},
			{
				Name:  "render",
				Usage: "uses Go templates to write static content",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Usage:   "path to the current project's `.stash.yml` configuration",
						EnvVars: []string{"CONFIG"},
					},
					&cli.BoolFlag{
						Name:  "watch",
						Usage: "render your templates as you make changes within your root directory",
					},
				},
				Action: func(c *cli.Context) error {
					return config.WrapWithConfig(c, commands.RenderHandler)
				},
			},
			{
				Name:  "server",
				Usage: "start a local web server to expose your rendered site",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Usage: "opens a browser window directly to the specified URL path",
					},
					&cli.BoolFlag{
						Name:  "no-browser",
						Usage: "prevent opening a new browser window on startup",
					},
				},
				Action: func(c *cli.Context) error {
					return config.WrapWithConfig(c, commands.ServerHandler)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
