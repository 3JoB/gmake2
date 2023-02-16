package main

import "github.com/urfave/cli/v2"

var CliFlag = []cli.Flag{
	CliFlagConfig,
	CliFlagDebug,
	CliFlagUpgrade,
	CliFlagUpgradeX,
}

var CliCommands = []*cli.Command{
	CliCommandInit,
	CliCommandUpdate,
	CliCommandVersion,
}

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

var CliFlagUpgradeX = &cli.BoolFlag{
	Name:  "x",
	Value: false,
	Usage: "Force an update to be downloaded from the server",
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

var VersionInfo = `GMake2 is distributed under Mozilla Public License 2.0.
Github: https://github.com/3JoB/gmake2

Version: ` + SoftVersion + ` (` + SoftVersionCode + `) [ Built on: ` + SoftBuildTime + ` ]
CommitID: ` + SoftCommit

var InitFileContent = `config:
  default: all
  proxy:
  req: false

var:
  msg: GMake2
  
all: |
  @echo Hello! {{.msg}}!`

var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52 GMake2/" + SoftVersion