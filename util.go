package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/goutil/fsutil"
)

func if_func(f []string, ym map[string]any) error {
	if f[3] != "then" {
		ErrPrintf("GMake2: Invalid operator at %v \n", f[3])
	}
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
	if f[5] != "or" {
		ErrPrintf("GMake2: Invalid operator at %v \n", f[5])
	}
	if f[6] == "null" {
		return nil
	}
	run(ym, f[6])
	return nil
}

/*func write(path string, v any) {
	checkError(fsutil.WriteFile(path, v, 0664))
}*/

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
		ErrPrintf("GMake2:  Something went wrong, you must examine the following error messages to determine what went wrong. \n %v \n", err)
	}
}

func Println(a ...any) {
	if debug {
		log.Fatal(a...)
	}
	fmt.Println(a...)
}

func ErrPrint(a ...any) {
	if debug {
		log.Panic(a...)
	}
	fmt.Println(a...)
	os.Exit(0)
}

func ErrPrintf(format string, v ...any) {
	if debug {
		log.Panicf(format, v...)
	}
	fmt.Printf(format, v...)
	os.Exit(0)
}
