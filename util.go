package main

import (
	"fmt"

	"github.com/gookit/goutil/fsutil"
)

func checkThen(th string) {
	if th != "then" {
		fmt.Printf("GMake2: Invalid operator at %v \n", th)
		return
	}
}

func checkOr(th string) {
	if th != "or" {
		fmt.Printf("GMake2: Invalid operator at %v \n", th)
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
		fmt.Println(err)
		return
	}
}