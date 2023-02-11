package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	ufs "github.com/3JoB/ulib/fsutil"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
)

var (
	SoftVersion string
	SoftCommit  string
)

func init() {
	JsonData = make(map[string]string)
}

func main() {
	app := &cli.App{
		Name:  "GMake2",
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
				Value:   "GMakefile.yml",
				Usage:   "GMake2 Config File",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "GMake2 Version",
				Action: func(ctx *cli.Context) error {
					fmt.Println(VersionInfo)
					return nil
				},
			},
			{
				Name:   "init",
				Usage:  "Initialize in the current directory.",
				Action: InitFile,
			},
			{
				Name: "update",
				Usage: "Check for GMake2 updates (not applicable to distributions installed via choco,apt)",
				Action: CheckUpdate,
			},
		},
		Action: CMD,
	}

	err := app.Run(os.Args)
	checkError(err)
}

func CMD(c *cli.Context) error {
	// Parsing GMakefile
	ym := parseConfig(c.String("c"))

	// Read debug information
	debug = c.Bool("debug")

	// Parse Map
	parseMap(ym)

	if cfg["proxy"] != nil {
		u, err := url.Parse(cfg["proxy"].(string))
		checkError(err)
		Client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(u)},
		}
	} else {
		Client = http.DefaultClient
	}

	commands_args := ""

	// Check if the initialization command group exists
	if cfg["init"].(bool) {
		if ym["init"] != nil {
			run(ym, "init")
		}
	}

	if c.Args().Len() != 0 {
		commands_args = c.Args().First()
	} else {
		if cast.ToString(cfg["default"]) != "" {
			commands_args = cast.ToString(cfg["default"])
		} else {
			commands_args = "all"
		}
	}

	// Execution command group
	run(ym, commands_args)
	return nil
}

// Create a GMakefile
func InitFile(c *cli.Context) error {
	// If GMakefile exists, make it wait 12 seconds
	if isFile("GMakefile.yml") {
		Println("GMake2: Note! There are already GMakefile.yml files in the directory! Now you still have 12 seconds to prevent GMAKE2 from covering the file!")
		time.Sleep(time.Second * 12)
		rm("GMakefile.yml")
		Println("GMake2: File is being covered.")
	}

	// Then write to the file
	if err := ufs.File("GMakefile.yml").SetTrunc().Write(InitFileContent); err != nil {
		ErrPrintf("GMake2: Error! %v \n", err.Error())
	}
	Println("GMake2: GMakefile.yml file has been generated in the current directory.")
	return nil
}

// Check for updates
func CheckUpdate(c *cli.Context) error {
	return nil
}