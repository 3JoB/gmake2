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
	"github.com/urfave/cli/v2"
)

var (
	vars map[string]any
	cfg  map[string]any
	ctx  *cli.Context
	R    Req
)

func run(ym map[string]any, commands string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		// fmt.Println(pair[0])
		vars[pair[0]] = pair[1]
	}
	cmdDir := ""
	if cast.ToString(ym[commands]) == "" {
		EPrintf("GMake2: Command not found %v \n", commands)
	}
	k, v := commands, ym[commands]
	if k != "vars" && k != "config" {
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
					EPrintf("GMake2: Illegal instruction!\n GMake2: Error Command: %v \n", fmt.Sprint(cmdStrs[:]))
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
				case "@val":
					varg := args[2:]
					vcmd := exec.Command(args[1], varg...)
					if cmdDir != "" {
						vcmd.Dir = cmdDir
					}
					val(args, vcmd)
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
				case "@sleep":
					time.Sleep(time.Second * cast.ToDuration(args[0]))
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
				case "@req":
					R.Network(args...)
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
