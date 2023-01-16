package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/3JoB/telebot/pkg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/goutil/fsutil"
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
		fmt.Println("GMake: Invalid operator!")
	}
	return nil
}

func checkThen(th string) {
	if th != "then" {
		fmt.Printf("GMake: Invalid operator at %v \n", th)
		return
	}
}

func checkOr(th string) {
	if th != "or" {
		fmt.Printf("GMake: Invalid operator at %v \n", th)
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

// TODO:
func get_json(r ...string) error {
	return nil
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
	if len(r) != 5 {
		fmt.Println("GMake: Illegal instruction!!!")
		os.Exit(0)
	}
	if _, err := url.Parse(r[1]); err != nil {
		fmt.Println("GMake: Url check failed!!!")
		fmt.Println("GMake: " + err.Error())
		os.Exit(0)
	}
	client := resty.New()
	resp, err := client.R().SetHeader("User-Agent", "github.com/3JoB/gmake2 grab/3").Get(r[1])
	checkError(err)
	if resp.StatusCode() != 200 {
		fmt.Printf("GMake: Server returned status code: %v \n", resp.StatusCode())
		os.Exit(0)
	}
	fmt.Printf("Parsing json from %v", r[1])
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
	fmt.Printf("GMake: Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("GMake: Connection info: %v\n", resp.HTTPResponse.Status)
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
			fmt.Printf("GMake: transferred %v/%v bytes (%.2f%%)\n", resp.BytesComplete(), fsize, 100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "GMake: Download failed: %v\n", err)
		return
	}

	fmt.Printf("GMake: Download saved to ./%v \n", resp.Filename)
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
