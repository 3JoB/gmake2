package main

import (
	"fmt"
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
		Name: "GMake2",
		Usage: "program like make",
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
				Aliases: []string{"v"},
				Usage: "GMake2 Version",
				Action: func(ctx *cli.Context) error {
					fmt.Println(VersionInfo)
					return nil
				},
			},
			{
				Name: "init",
				Usage: "Initialize in the current directory.",
				Action: InitFile,
			},
		},
		Action: CMD,
	}

	err := app.Run(os.Args)
	checkError(err)
}

func CMD(c *cli.Context) error {
	ctx = c
	ym := parseConfig(c.String("c"))
	parseMap(ym)
	commands_args := ""
	if len(c.Args().Slice()) != 0 {
		commands_args = c.Args().First()
	} else {
		if cast.ToString(cfg["default"]) != "" {
			commands_args = cast.ToString(cfg["default"])
		} else {
			commands_args = "all"
		}
	}
	run(ym, commands_args)
	return nil
}