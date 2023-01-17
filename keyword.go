package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/3JoB/telebot/pkg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

func ifelse(ym map[string]any, f []string) error {
	switch f[1] {
	case "==":
		if f[0] == f[2] {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	case "!=":
		if f[0] != f[2] {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	case "<":
		if cast.ToInt64(f[0]) < cast.ToInt64(f[2]) {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	case "<=":
		if cast.ToInt64(f[0]) <= cast.ToInt64(f[2]) {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	case ">":
		if cast.ToInt64(f[0]) > cast.ToInt64(f[2]) {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	case ">=":
		if cast.ToInt64(f[0]) >= cast.ToInt64(f[2]) {
			return if_func(f, ym)
		}
		return if_func2(f, ym)
	default:
		EPrint("GMake2: Invalid operator!")
	}
	return nil
}

func val(r []string, c *exec.Cmd) {
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	checkError(err)
	outStr, errStr := pkg.String(stdout.Bytes()), pkg.String(stderr.Bytes())
	if errStr != "" {
		EPrint(errStr)
	}
	vars[r[0]] = outStr
}

/*
api.json:

	{
		"msg": "666"
	}

@json url https://example.com/api.json string msg vb

@echo {{.vb}}
*/
func get_json_url(r []string) error {
	fmt.Println(r[:])
	if len(r) != 5 {
		EPrint("GMake2: Illegal instruction!!!")
	}
	if _, err := url.Parse(r[1]); err != nil {
		fmt.Println("GMake2: Url check failed!!!")
		EPrint("GMake2: " + err.Error())
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52").
		SetHeader("APP-User-Agent", "github.com/3JoB/gmake2 grab/3").
		Get(r[1])
	checkError(err)
	if resp.StatusCode() != 200 {
		EPrintf("GMake2: Server returned status code: %v \n", resp.StatusCode())
	}
	fmt.Printf("Parsing json from %v \n", r[1])
	result := gjson.Get(pkg.String(resp.Body()), r[3])
	switch r[2] {
	case "string", "String":
		vars[r[4]] = result.String()
	case "bool", "Bool":
		vars[r[4]] = result.Bool()
	case "int", "int8", "int16", "int32", "int64":
		vars[r[4]] = result.Int()
	case "uint", "utin8", "uint16", "uint32", "uint64":
		vars[r[4]] = result.Uint()
	case "float", "float32", "float64":
		vars[r[4]] = result.Float()
	default:
		vars[r[4]] = result.String()
	}
	return nil
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
	fmt.Printf("GMake2: Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("GMake2: Connection info: %v\n", resp.HTTPResponse.Status)
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
			fmt.Printf("GMake2: transferred %v/%v bytes (%.2f%%)\n", resp.BytesComplete(), fsize, 100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "GMake2: Download failed: %v\n", err)
		return
	}

	fmt.Printf("GMake2: Download saved to ./%v \n", resp.Filename)
}

func copy(src, dst string) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if isDir(src) {
		if !isDir(dst) {
			EPrintf("gmake2: Cannot copy directory to file src=%v dst=%v", src, dst)
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
