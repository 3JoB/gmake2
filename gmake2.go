package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
)

var (
	vars     map[string]any
	cfg      map[string]any
	JsonData map[string]string
	R        Req
	Client   *http.Client
	debug    bool
)

func init() {
	JsonData = make(map[string]string)
}

func run(ym map[string]any, commands string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		// fmt.Println(pair[0])
		vars[pair[0]] = pair[1]
	}
	if ym[commands] == nil {
		ErrPrintf("GMake2: Command not found %v \n", commands)
	}
	cmdDir := ""
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
					ErrPrintf("GMake2: Illegal instruction!\nGMake2: Error Command: %v \n", fmt.Sprint(cmdStrs[:]))
				}
				switch bin {
				case "@var":
					vars[args[0]] = strings.Join(args[1:], " ")
				case "@env":
					os.Setenv(args[0], strings.Join(args[1:], " "))
				case "@cmd":
					run(ym, args[0])
				case "@wait":
					wait(args...)
				case "@if":
					ifelse(ym, args)
				case "@val":
					arg := args[2:]
					cmd := exec.Command(args[1], arg...)
					if cmdDir != "" {
						cmd.Dir = cmdDir
					}
					val(args, cmd)
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
					JsonUrl(args)
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
					if cast.ToBool(cfg["req"]) {
						R.Network(args...)
					} else {
						Println("GMake2: The @req tag has been deprecated.")
					}
				case "@async":
					go run(ym, args[0])
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
	Println(c.String())
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
