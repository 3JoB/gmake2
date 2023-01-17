package main

import (
	"os"

	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
)

var (
	SoftVersion string
	SoftCommit string
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "debug mode",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "gmake2.yml",
				Usage:   "GMake2 Config File",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "version",

			},
		},
		Action: func(c *cli.Context) error {
			ctx = c
			ym := parseConfig(cfgFile)
			parseMap(ym)
			commands_args := ""
			if c.NArg() != 1 {
				if cfg["default"] != "" {
					commands_args = cast.ToString(cfg["default"])
				} else {
					commands_args = "all"
				}
			} else {
				commands_args = c.Args().First()
			}
			run(ym, commands_args)
			return nil
		},
	}

	err := app.Run(os.Args)
	checkError(err)
}
