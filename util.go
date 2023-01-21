package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gookit/goutil/fsutil"
	"github.com/urfave/cli/v2"
)

func checkThen(th string) {
	if th != "then" {
		EPrintf("GMake2: Invalid operator at %v \n", th)
		return
	}
}

func checkOr(th string) {
	if th != "or" {
		EPrintf("GMake2: Invalid operator at %v \n", th)
		return
	}
}

func if_func(f []string, ym map[string]any) error {
	checkThen(f[3])
	if f[4] == "null" {
		return nil
	}
	run(ym, f[4])
	return nil
}

func if_func2(f []string, ym map[string]any) error {
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

func write(path string, v any) {
	checkError(fsutil.WriteFile(path, v, 0664))
}

func copyFile(src, dst string) error {
	return fsutil.CopyFile(src, dst)
}

func isDir(path string) bool {
	return fsutil.IsDir(path)
}

func isFile(path string) bool {
	return fsutil.IsFile(path)
}

func checkError(err error) {
	if err != nil {
		EPrint(err)
	}
}

func EPrint(a ...any) {
	if ctx.Bool("debug") {
		log.Fatal(a...)
	}
	fmt.Println(a...)
	os.Exit(0)
}

func EPrintf(format string, v ...any) {
	if ctx.Bool("debug") {
		panic(fmt.Sprintf(format, v...))
	}
	fmt.Printf(format, v...)
	os.Exit(0)
}

func InitFile(c *cli.Context) error {
	if isFile("GMakefile.yml") {
		fmt.Println("GMake2: Note! There are already GMakefile.yml files in the directory! Now you still have 12 seconds to prevent GMAKE2 from covering the file!")
		time.Sleep(time.Second * 12)
		rm("GMakefile.yml")
		fmt.Println("GMake2: File is being covered.")
	}
	touch("GMakefile.yml")
	write("GMakefile.yml", InitFileContent)
	fmt.Println("GMake2: GMakefile.yml file has been generated in the current directory.")
	return nil
}
