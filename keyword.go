package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/3JoB/ulib/json"
	"github.com/3JoB/ulib/reflect"
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
		ErrPrintf("GMake2: Invalid operator!\nGMake2: Error Command: %v \n", fmt.Sprint(f[:]))
	}
	return nil
}

func val(r []string, c *exec.Cmd) {
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	checkError(err)
	outStr, errStr := reflect.String(stdout.Bytes()), reflect.String(stderr.Bytes())
	if errStr != "" {
		ErrPrintf("GMake: Val Failed!!!\nGMake2: Error Command: %v \n", errStr)
	}
	vars[r[0]] = outStr
}

func JsonUrl(r []string) error {
	switch r[0] {
	case "parse":
		if len(r) != 4 {
			ErrPrintf("GMake2: Illegal instruction!!!\nGMake2: Error Command: %v \n", fmt.Sprint(r[:]))
		}
		vars[r[3]] = gjson.Get(JsonData[r[1]], r[2]).String()
	default:
		if len(r) != 2 {
			ErrPrintf("GMake2: Illegal instruction!!!\nGMake2: Error Command: %v \n", fmt.Sprint(r[:]))
		}
		if _, err := url.Parse(r[0]); err != nil {
			ErrPrint("GMake2: Url check failed!!!\nGMake2: " + err.Error())
		}

		resp, err := resty.NewWithClient(Client).R().
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 GMake2/"+SoftVersion).
			Get(r[0])

		checkError(err)

		if resp.StatusCode() != 200 {
			ErrPrintf("GMake2: Server returned status code: %v \n", resp.StatusCode())
		}

		rd := reflect.String(resp.Body())

		if rd != "" {
			JsonData[r[1]] = rd
		}
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
	f, err := fsutil.CreateFile(path, 0664, 0666)
	checkError(err)
	f.Close()
}

func downloadFile(filepath string, url string) {
	client := grab.NewClient()
	client.HTTPClient = Client
	client.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 GMake2/" + SoftVersion

	req, _ := grab.NewRequest(filepath, url)

	fmt.Printf("GMake2: Downloading %v...\n", req.URL())

	resp := client.Do(req)

	fmt.Printf("GMake2: Connection info: %v\n", resp.HTTPResponse.Status)

	fsize := cast.ToString(resp.Size)

	if fsize == "" {
		fsize = "unknown"
	}

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
			ErrPrintf("GMake2: Cannot copy directory to file src=%v dst=%v \n", src, dst)
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

type Req struct {
	Header map[string]string
	Body   any
	File   string
	Method string
	Uri    string
	Value  string
	Req    *resty.Request
	Resp   *resty.Response
}

// Deprecated: This method is about to be deprecated and is no longer supported
func (r *Req) Network(str ...string) {
	switch str[0] {
	case "def":
		d := strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		switch str[1] {
		case "header":
			headers := make(map[string]string)
			json.UnmarshalString(d, &headers)
			r.Header = headers
		case "body":
			r.Body = d
		case "file":
			r.File = d
		case "method":
			r.Method = d
		case "uri", "url":
			r.Uri = d
		case "value":
			r.Value = d
		default:
			Println("GMake2: @req: unknown method: " + fmt.Sprint(str[1:]))
		}
	default:
		r.Request()
	}
	// v := strings.ReplaceAll(strings.Trim(fmt.Sprint(str), "[]"), " ", " ")
}

// Deprecated
//
// This method is about to be deprecated and is no longer supported
func (r *Req) Request() {
	_, err := url.Parse(r.Uri)
	checkError(err)

	client := resty.NewWithClient(Client)

	if r.Header == nil {
		r.Header = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52 GMake2/" + SoftVersion,
		}
	}
	if r.Header["User-Agent"] == "" {
		r.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52 GMake2/" + SoftVersion
	}

	r.Req = client.R().SetHeaders(r.Header)

	if r.Body != nil {
		r.Req = r.Req.SetBody(r.Body)
	}

	if r.File != "" {
		r.Req = r.Req.SetFile(r.File, r.File)
	}

	switch r.Method {
	case "POST", "post":
		r.Resp, err = r.Req.Post(r.Uri)
	case "DELETE", "delete":
		r.Resp, err = r.Req.Delete(r.Uri)
	case "PATCH", "patch":
		r.Resp, err = r.Req.Patch(r.Uri)
	case "PUT", "put":
		r.Resp, err = r.Req.Put(r.Uri)
	default:
		r.Resp, err = r.Req.Get(r.Uri)
	}

	checkError(err)

	defer r.Resp.RawBody().Close()

	if r.Resp.StatusCode() != 200 {
		Println("GMake2: @req: Server returned error code:" + cast.ToString(r.Resp.StatusCode()))
	} else {
		Println("GMake2: @req: 200 ok")
	}

	body := reflect.String(r.Resp.Body())
	if body != "" {
		if r.Value != "" {
			vars[r.Value] = body
		}
		Println(body)
	}
}

func wait(v ...string) {
	ar := len(v)
	if ar < 2 {
		ErrPrint("GMake2: @wait bad format!")
	}
	// v[0] print
	Println(v[:ar-1])
	var t string
	fmt.Print("=> ")
	fmt.Scan(&t)
	vars[v[ar-1]] = t
}
