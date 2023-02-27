package main

import (
	"os"

	"github.com/3JoB/ulib/fsutil"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "GMake2",
		Usage:    "Lightning-like GMake-like programs.",
		Before:   CliBeforeFunc,
		Flags:    CliFlag,
		Commands: CliCommands,
		Action:   CliAction,
	}

	err := app.Run(os.Args)
	checkError(err)
}

func CliBeforeFunc(c *cli.Context) error {
	// Read debug information
	debug = c.Bool("debug")
	Tags = c.String("tag")
	return nil
}

func CliAction(c *cli.Context) error {
	// Parsing GMakefile
	ym := parseConfig(c.String("c"))
	// Parse Map
	parseMap(ym)
	// Parse Tags
	if Tags != "" {
		parseTags(Tags)
	}
	// Import Proxy Config
	ImportProxy(cfg["proxy"])

	commands_args := ""

	// Check if the initialization command group exists
	if ym["init"] != nil {
		run(ym, "init")
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
	if fsutil.IsExist(e) {
		Println("GMake2: Note! There are already GMakefile.yml files in the directory! Now you still have 12 seconds to prevent GMAKE2 from covering the file!")
		sleep(12)
		remove(e)
		Println("GMake2: File is being covered.")
	}

	// Then write to the file
	checkError(write(e, InitFileContent))
	/*if err := ufs.File("GMakefile.yml").SetTrunc().Write(InitFileContent); err != nil {
		ErrPrintf("GMake2: Error! %v \n", err.Error())
	}*/
	Println("GMake2: GMakefile.yml file has been generated in the current directory.")
	return nil
}
