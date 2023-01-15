package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile string
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var rootCmd = &cobra.Command{
		Use:   "gmake",
		Short: "parse custom makefile and execute",
		Long:  "",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ym := parseConfig(cfgFile)
			cmds := ""
			if len(cmd.Flags().Args()) != 1 {
				cmds = "all"
			} else {
				cmds = cmd.Flags().Args()[0]
			}
			run(ym, cmds)
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
	t := time.Now().Format("2006-01-02 15:04")
	vars["time"] = t
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
				if err != nil {
					log.Fatal(err)
				}
				for i, cmdStr := range cmdStrs {
					cmdStrs[i] = ResolveVars(vars, cmdStr)
				}
				bin, args := cmdStrs[0], cmdStrs[1:]
				if len(args) != 1 {
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
					err := downloadFile(args[1], args[0])
					checkError(err)
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
	
	/*for k, v := range ym {
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
					if err != nil {
						log.Fatal(err)
					}
					for i, cmdStr := range cmdStrs {
						cmdStrs[i] = ResolveVars(vars, cmdStr)
					}
					bin, args := cmdStrs[0], cmdStrs[1:]
					switch bin {
					case "@var":
						vars[args[0]] = strings.Join(args[1:], " ")
					case "@env":
						os.Setenv(args[0], strings.Join(args[1:], " "))
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
						err := downloadFile(args[1], args[0])
						checkError(err)
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
	}*/
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
		return []string{}, errors.New(`Unclosed quote in command line: " ` + command + ` "`)
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
	err := os.RemoveAll(path)
	checkError(err)
}

func mkdir(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	checkError(err)
}

func touch(path string) {
	f, err := os.Create(path)
	checkError(err)
	f.Close()
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func copy(src, dst string) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if isDir(src) {
		if !isDir(dst) {
			panic(fmt.Errorf("不能复制目录到文件 src=%v dst=%v", src, dst))
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

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func ResolveVars(vars any, templateStr string) string {
	if vars == nil {
		return templateStr
	}
	t := template.Must(template.New("template").Parse(templateStr))
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, vars); err != nil {
		log.Fatal(err)
	}
	s := buf.String()
	return s
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false
	}
	return s == nil || s.IsDir()
}

func isFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false
	}
	return s == nil || !s.IsDir()
}

func ExecCmd(c *exec.Cmd) {
	log.Println(c.String())
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
