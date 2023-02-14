package main

import "github.com/urfave/cli/v2"

var CliFlagDebug = &cli.BoolFlag{
	Name:  "debug",
	Value: false,
	Usage: "debug mode",
}

var CliFlagConfig = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Value:   "GMakefile.yml",
	Usage:   "GMake2 Config File",
}

var CliFlagUpgrade = &cli.BoolFlag{
	Name:  "upgrade",
	Value: false,
	Usage: "Mandatory upgrade channel edition",
}

var CliCommandVersion = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "GMake2 Version",
	Action: func(ctx *cli.Context) error {
		Println(VersionInfo)
		return nil
	},
}

var CliCommandInit = &cli.Command{
	Name:   "init",
	Usage:  "Initialize in the current directory.",
	Action: InitFile,
}

var CliCommandUpdate = &cli.Command{
	Name:   "update",
	Usage:  "Check for GMake2 updates (not applicable to distributions installed via choco,apt)",
	Action: CheckUpdate,
}
