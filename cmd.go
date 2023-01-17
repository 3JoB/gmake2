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

var (
	VersionInfo string = `
GMake2 is distributed under Apache-2.0 license.
Github: https://github.com/3JoB/gmake2

Version: `+SoftVersion+`
CommitID: `+SoftCommit
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
				Aliases: []string{"v"},
				Usage: "GMake2 Version",
				Action: func(ctx *cli.Context) error {
					fmt.Println(VersionInfo)
					return nil
				},
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
			fmt.Printf("这是cfg: %v \n", commands_args)
		} else {
			commands_args = "all"
			fmt.Printf("普通cfg: %v \n", commands_args)
		}
	}
	run(ym, commands_args)
	return nil
}