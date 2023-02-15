package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
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

func request(url string) *resty.Response {
	if Client == nil {
		Client = http.DefaultClient
	}
	resp, err := resty.NewWithClient(Client).R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 GMake2/"+SoftVersion).
		Get(url)
	checkError(err)
	defer resp.RawBody().Close()

	return resp
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
		ErrPrintf("GMake2:  Something went wrong, you must examine the following error messages to determine what went wrong. \n%v \n", err)
	}
}

func Println(a ...any) {
	if debug {
		log.Fatal(a...)
	} else {
		fmt.Println(a...)
	}
}

func Printf(format string, v ...any) {
	if debug {
		log.Fatalf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
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
