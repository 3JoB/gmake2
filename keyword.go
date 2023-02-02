package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	//"strings"
	"time"

	"github.com/3JoB/telebot/pkg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	//"github.com/goccy/go-json"
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
		EPrintf("GMake2: Invalid operator!\n GMake2: Error Command: %v \n", fmt.Sprint(f[:]))
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
		EPrintf("GMake: Val Failed!!!\n  GMake2: Error Command: %v \n", errStr)
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
	if len(r) != 5 {
		EPrintf("GMake2: Illegal instruction!!!\nGMake2: Error Command: %v \n", fmt.Sprint(r[:]))
	}
	if _, err := url.Parse(r[1]); err != nil {
		EPrint("GMake2: Url check failed!!!\nGMake2: " + err.Error())
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.52").
		SetHeader("APP-User-Agent", "github.com/3JoB/gmake2 Version/2").
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
	case "uint", "uint8", "uint16", "uint32", "uint64":
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
	f, err := fsutil.CreateFile(path, 0664, 0666)
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
			EPrintf("GMake2: Cannot copy directory to file src=%v dst=%v \n", src, dst)
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

/*type Req struct {
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
		switch str[1] {
		case "header":
			headers := make(map[string]string)
			json.Unmarshal(pkg.Bytes(strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")), &headers)
			r.Header = headers
		case "body":
			r.Body = strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		case "file":
			r.File = strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		case "method":
			r.Method = strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		case "uri", "url":
			r.Uri = strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		case "value":
			r.Value = strings.ReplaceAll(strings.Trim(fmt.Sprint(str[2:]), "[]"), " ", " ")
		default:
			fmt.Println("GMake2: @req: unknown method: " + fmt.Sprint(str[1:]))
		}
	default:
		r.Request()
	}
	// v := strings.ReplaceAll(strings.Trim(fmt.Sprint(str), "[]"), " ", " ")
}

// Deprecated
//
//This method is about to be deprecated and is no longer supported
func (r *Req) Request() {
	_, err := url.Parse(r.Uri)
	checkError(err)
	client := resty.New()
	r.Req = client.R().SetHeaders(r.Header).SetBody(r.Body)
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
		fmt.Println("GMake2: @req: Server returned error code:" + cast.ToString(r.Resp.StatusCode()))
	} else {
		fmt.Println("GMake2: @req: 200 ok")
	}
	body := pkg.String(r.Resp.Body())
	if body != "" {
		if r.Value != "" {
			vars[r.Value] = body
		}
		fmt.Println(body)
	}
}
*/