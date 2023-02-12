package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	ufs "github.com/3JoB/ulib/fsutil"
	"github.com/3JoB/ulib/reflect"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
)

var (
	SoftVersion     string
	SoftVersionCode string
	SoftCommit      string
)

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
			&cli.BoolFlag{
				Name: "upgrade",
				Value: false,
				Usage: "Mandatory upgrade channel edition",
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
				Name:   "update",
				Usage:  "Check for GMake2 updates (not applicable to distributions installed via choco,apt)",
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
		u, err := url.Parse(cast.ToString(cfg["proxy"]))
		checkError(err)
		Client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(u)},
		}
	} else {
		Client = http.DefaultClient
	}

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
	run_path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	downloadPath := ""
	resp, err := resty.New().R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52").
		SetHeader("APP-User-Agent", "github.com/3JoB/gmake2 Version/2").
		Get("https://lcag.org/gmake2.raw")

	checkError(err)
	if resp.StatusCode() != 200 {
		ErrPrintf("GMake2: Server returned status code: %v \n", resp.StatusCode())
	}

	defer resp.RawBody().Close()

	rd := reflect.String(resp.Body())

	version_code := gjson.Get(rd, "version_code").Int()
	version := gjson.Get(rd, "version").String()
	update_url := gjson.Get(rd, "url").String()

	if version_code > cast.ToInt64(SoftVersionCode) {
		if c.Bool("upgrade") {
			run_path = "n"
		}
		switch run_path {
		case `C:\ProgramData\chocolatey\lib\gmake2\tools`:
			Println("Sorry, Chocolatey does not support automatic updates, please use the command 'choco update gmake2 --version="+version+"' to update gmake2")
		case "/usr/bin":
			Println("Sorry, apt does not support automatic updates, please use the command 'apt update && apt upgrade' to update gmake2")
		default:
			filename:= "gmake2"
			if cast.ToString(vars["runtime.os"]) == "windows" {
				filename = filename + ".exe"
				downloadPath = run_path + `\` + filename
			} else {
				downloadPath = run_path + `/` + filename
			}
			downloadUrl := update_url+"?arch="+cast.ToString(vars["runtime.arch"])+"&os="+cast.ToString(vars["runtime.os"])+"&version="+version

			downloadFile(downloadPath, downloadUrl)

			Println("GMake2 has been updated to "+version+"("+cast.ToString(version_code)+")")
		}
	}
	return nil
}
