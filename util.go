package main

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	ufs "github.com/3JoB/ulib/fsutil"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

func operation_1(f []string, ym map[string]any) error {
	if f[3] != "then" {
		ErrPrintf("GMake2: Invalid operator at %v \n", f[3])
	}
	if f[4] == "null" {
		return nil
	}
	run(ym, f[4])
	return nil
}

func operation_2(f []string, ym map[string]any) error {
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

func request(url string) *resty.Response {
	if Client == nil {
		Client = http.DefaultClient
	}
	resp, err := resty.NewWithClient(Client).R().
		SetHeader("User-Agent", UserAgent).
		Get(url)
	checkError(err)
	return resp
}

func write(path, v string) error {
	return ufs.File(path).SetTrunc().Write(v)
}

func copyFile(src, dst string) error {
	return ufs.File(src).CopyTo(dst)
}

func isDir(path string) bool {
	return ufs.IsDir(path)
}

func isFile(path string) bool {
	return ufs.IsFile(path)
}

func sleep(t any) {
	time.Sleep(time.Second * cast.ToDuration(t))
}

func replace(v []string) string {
	return strings.ReplaceAll(strings.Trim(fmt.Sprint(v), "[]"), " ", " ")
}

func checkError(err error) {
	if err != nil {
		ErrPrintf("GMake2: Something went wrong, you must examine the following error messages to determine what went wrong. \n%v \n", err)
	}
}

func Println(a ...any) {
	fmt.Println(a...)
}

func Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

func ErrPrint(a ...any) {
	if debug {
		panic(fmt.Sprintln(a...))
	}
	fmt.Println(a...)
	Exit()
}

func ErrPrintf(format string, v ...any) {
	if debug {
		panic(fmt.Sprintf(format, v...))
	}
	fmt.Printf(format, v...)
	Exit()
}

func Exit() {
	os.Exit(0)
}

// from: https://github.com/cavaliergopher/grab/v3
func guessFilename(resp *http.Response) (string, error) {
	filename := resp.Request.URL.Path
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			if val, ok := params["filename"]; ok {
				filename = val
			} // else filename directive is missing.. fallback to URL.Path
		}
	}

	// sanitize
	if filename == "" || strings.HasSuffix(filename, "/") || strings.Contains(filename, "\x00") {
		return "GMakeDL.tmp", nil
	}

	filename = filepath.Base(path.Clean("/" + filename))
	if filename == "" || filename == "." || filename == "/" {
		return "GMakeDL.tmp", nil
	}

	return filename, nil
}

func ImportProxy(v any) {
	if v != nil {
		u, err := url.Parse(cast.ToString(cfg["proxy"]))
		checkError(err)
		Client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(u)},
		}
	} else {
		Client = http.DefaultClient
	}
}
