package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gookit/goutil/fsutil"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile       string
	commands_args string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gmake",
		Short: "parse custom makefile and execute",
		Long:  "",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ym := parseConfig(cfgFile)
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
	var vars map[string]any
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
		fmt.Printf("gmake: Command not found %v \n", commands)
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
					fmt.Println("gmake: Illegal instruction!")
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
	v["time"] = time.Now().Format("2006-01-02 15:04")
	v["time_utc"] = time.Now().UTC().Format("2006-01-02 15:04")
	v["time_unix"] = time.Now().Unix()
	v["time_utc_unix"] = time.Now().UTC().Unix()
	v["runtime_os"] = runtime.GOOS
	v["runtime_arch"] = runtime.GOARCH
	return v
}

func ifelse(ym map[string]any, f []string) error {
	switch f[1] {
	case "==":
		if f[0] == f[2] {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	case "!=":
		if f[0] != f[2] {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	case "<":
		if cast.ToInt64(f[0]) < cast.ToInt64(f[2]) {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	case "<=":
		if cast.ToInt64(f[0]) <= cast.ToInt64(f[2]) {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	case ">":
		if cast.ToInt64(f[0]) > cast.ToInt64(f[2]) {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	case ">=":
		if cast.ToInt64(f[0]) >= cast.ToInt64(f[2]) {
			return ifunc(f, ym)
		}
		return ifunc2(f, ym)
	default:
		fmt.Println("gmake: Invalid operator!")
	}
	return nil
}

func checkThen(th string) {
	if th != "then" {
		fmt.Printf("gmake: Invalid operator at %v \n", th)
		return
	}
}

func checkOr(th string) {
	if th != "or" {
		fmt.Printf("gmake: Invalid operator at %v \n", th)
		return
	}
}

func ifunc(f []string, ym map[string]any) error {
	checkThen(f[3])
	if f[4] == "null" {
		return nil
	}
	run(ym, f[4])
	return nil
}

func ifunc2(f []string, ym map[string]any) error {
	if len(f) != 7 {
		return nil
	}
	checkOr(f[5])
	if f[6] == "null" {
		return nil
	}
	run(ym, f[6])
	return nil
}

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(`gmake: Unclosed quote in command line: " ` + command + ` "`)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

func mv(from, to string) {
	copy(from, to)
	rm(from)
}

func rm(path string) {
	checkError(os.RemoveAll(path))
}

func mkdir(path string) {
	checkError(os.MkdirAll(path, os.ModePerm))
}

func touch(path string) {
	f, err := os.Create(path)
	checkError(err)
	f.Close()
}

func downloadFile(filepath string, url string) {
	// Get the data
	client := grab.NewClient()
	client.UserAgent = "github.com/3JoB/gmake2 grab/3"
	req, _ := grab.NewRequest(filepath, url)
	// start download
	fmt.Printf("gmake: Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("gmake: Connection info: %v\n", resp.HTTPResponse.Status)
	fsize := cast.ToString(resp.Size)
	if fsize == "" {
		fsize = "unknown"
	}

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("gmake: transferred %v/%v bytes (%.2f%%)\n", resp.BytesComplete(), fsize, 100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "gmake: Download failed: %v\n", err)
		return
	}

	fmt.Printf("gmake: Download saved to ./%v \n", resp.Filename)
}

func copy(src, dst string) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if isDir(src) {
		if !isDir(dst) {
			fmt.Printf("gmake2: Cannot copy directory to file src=%v dst=%v", src, dst)
			return
		}
		si, err := os.Stat(src)
		checkError(err)
		// dst = path.Join(dst, filepath.Base(src))
		err = os.MkdirAll(dst, si.Mode())
		checkError(err)
		entries, err := os.ReadDir(src)
		checkError(err)
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			if entry.IsDir() {
				copy(srcPath, dstPath)
			} else {
				// Skip symlinks.
				if entry.Type()&os.ModeSymlink != 0 {
					continue
				}
				copyFile(srcPath, dstPath)
			}
		}
	} else {
		if isFile(dst) {
			copyFile(src, dst)
		} else {
			copyFile(src, path.Join(dst, filepath.Base(src)))
		}
	}
}

func copyFile(src, dst string) error {
	return fsutil.CopyFile(src, dst)
}

func ResolveVars(vars any, templateStr string) string {
	if vars == nil {
		return templateStr
	}
	t := template.Must(template.New("template").Parse(templateStr))
	buf := new(bytes.Buffer)
	checkError(t.Execute(buf, vars))
	s := buf.String()
	return s
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func isDir(path string) bool {
	return fsutil.IsDir(path)
}

func isFile(path string) bool {
	return fsutil.IsFile(path)
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
