package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile string
	vars    map[string]any
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gmake",
		Short: "parse custom makefile and execute",
		Long:  "",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ym := parseConfig(cfgFile)
			commands_args := ""
			if len(cmd.Flags().Args()) != 1 {
				commands_args = "all"
			} else {
				commands_args = cmd.Flags().Args()[0]
			}
			run(ym, commands_args)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "gmake.yml", "config file")
	rootCmd.Execute()
}

func parseConfig(cfgFile string) map[string]any {
	ymlData, err := os.ReadFile(cfgFile)
	checkError(err)
	m := make(map[string]any)
	err = yaml.Unmarshal(ymlData, &m)
	checkError(err)
	return m
}

func run(ym map[string]any, commands string) {
	if v, ok := ym["vars"]; ok {
		vars = v.(map[string]any)
	} else {
		vars = make(map[string]any)
	}
	vars = variable(vars)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		// fmt.Println(pair[0])
		vars[pair[0]] = pair[1]
	}
	cmdDir := ""
	if cast.ToString(ym[commands]) == "" {
		fmt.Printf("GMake2: Command not found %v \n", commands)
		return
	}
	k, v := commands, ym[commands]
	if k != "vars" {
		lines := strings.Split(cast.ToString(v), "\n")
		for _, line := range lines {
			if line != "" {
				// 注释
				if strings.TrimSpace(line)[0] == '#' {
					continue
				}
				// line = ResolveVars(vars, line)
				cmdStrs, err := shellquote.Split(line)
				checkError(err)
				for i, cmdStr := range cmdStrs {
					cmdStrs[i] = ResolveVars(vars, cmdStr)
				}
				bin, args := cmdStrs[0], cmdStrs[1:]
				if len(args) == 0 {
					fmt.Println("GMake2: Illegal instruction!")
					return
				}
				switch bin {
				case "@var":
					vars[args[0]] = strings.Join(args[1:], " ")
				case "@env":
					os.Setenv(args[0], strings.Join(args[1:], " "))
				case "@cmd":
					run(ym, args[0])
				case "@if":
					ifelse(ym, args)
				case "#":
				case "@echo":
					fmt.Println(strings.Join(args, " "))
				case "@mv":
					mv(args[0], args[1])
				case "@copy":
					copy(args[0], args[1])
				case "@rm":
					rm(args[0])
				case "@json":
					get_json_url(args)
				case "@mkdir":
					mkdir(args[0])
				case "@touch":
					touch(args[0])
				case "@download":
					if len(args) == 1 {
						downloadFile(".", args[0])
					} else {
						downloadFile(args[1], args[0])
					}
				case "@cd":
					abs, err := filepath.Abs(args[0])
					checkError(err)
					cmdDir = abs
				default:
					cmd := exec.Command(bin, args...)
					if cmdDir != "" {
						cmd.Dir = cmdDir
					}
					ExecCmd(cmd)
				}
			}
		}
	}
}

func variable(v map[string]any) map[string]any {
	v["time"] = map[string]any{
		"now":      time.Now().Format("2006-01-02 15:04"),
		"utc":      time.Now().UTC().Format("2006-01-02 15:04"),
		"unix":     time.Now().Unix(),
		"utc_unix": time.Now().UTC().Unix(),
	}
	v["runtime"] = map[string]any{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
	return v
}

func ExecCmd(c *exec.Cmd) {
	fmt.Println(c.String())
	stdout, err := c.StdoutPipe()
	checkError(err)
	stderr, err := c.StderrPipe()
	checkError(err)
	err = c.Start()
	checkError(err)
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stderr, stderr)
	c.Wait()
}
