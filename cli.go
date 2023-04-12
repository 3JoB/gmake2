package main

import (
	"github.com/urfave/cli/v2"
)

var (
	SoftVersion     string
	SoftVersionCode string
	SoftBuildTime   string
	SoftCommit      string
	Tags            string

	// Config file (GMakefile)
	e string = "G.toml"
)

var VersionInfo = Sprintf(`GMake2 is distributed under Mozilla Public License 2.0.
Github: https://github.com/3JoB/gmake2

Version: %v (%v)
Built on: %v
CommitID: %v`, SoftVersion, SoftVersionCode, SoftBuildTime, SoftCommit)

var InitFileContent = `config:
  default: all
  proxy:
  req: false

vars:
  msg: GMake2
  
all: |
  @echo Hello! {{.msg}}!`

// User-Agent
var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52 GMake2/" + SoftVersion

var CliFlag = []cli.Flag{
	CliFlagConfig,
	CliFlagTags,
	CliFlagDebug,
}

var CliCommands = []*cli.Command{
	CliCommandInit,
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
	Value:   e,
	Usage:   "gmake2 Config File",
}

var CliFlagTags = &cli.StringFlag{
	Name:    "tag",
	Hidden:  true,
	Aliases: []string{"t"},
	Value:   Tags,
	Usage:   "gmake2 Config File",
}

var CliCommandVersion = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "gmake2 Version",
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

var Error_Invalid error = Errors("invalid operator")
