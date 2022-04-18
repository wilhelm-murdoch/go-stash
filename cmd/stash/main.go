package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/wilhelm-murdoch/go-stash/cmd/stash/actions"

	"github.com/urfave/cli/v2"
)

var (
	// username = "atapas"
	// hostname = "blog.greenroots.info"
	username = "BlitzkriegPunk"
	hostname = "wilhelm.codes"
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Usage:   "the @username of the target Hashnode user.",
				Value:   "",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "hostname",
				Usage:   "the hostname of the target Hashnode blog.",
				Value:   "",
				Aliases: []string{"h"},
			},
			&cli.StringFlag{
				Name:    "since",
				Usage:   ".",
				Value:   "",
				Aliases: []string{"s"},
			},
		},
		Action: actions.RootHandler,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
