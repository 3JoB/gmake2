package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/goutil/fsutil"
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
	if ctx.Bool("debug"){
		log.Fatal(a...)
	}
	fmt.Println(a...)
	os.Exit(0)
}

func EPrintf(format string, v ...any){
	if ctx.Bool("debug"){
	log.Fatalf(format, v...)
	}
	fmt.Printf(format, v...)
	os.Exit(0)
}