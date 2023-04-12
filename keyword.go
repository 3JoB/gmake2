package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"

	"github.com/3JoB/ulib/fsutil"
	"github.com/3JoB/ulib/json"
	"github.com/3JoB/unsafeConvert"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

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

func ExecCmd(c *exec.Cmd) error {
	Println(c.String())
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stderr, stderr)
	c.Wait()
	return nil
}

func val(r []string, c *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	if err := c.Run(); err != nil {
		return err
	}
	outStr, errStr := unsafeConvert.StringReflect(stdout.Bytes()), unsafeConvert.StringReflect(stderr.Bytes())
	if errStr != "" {
		return Errors(errStr)
	}
	vars[r[0]] = outStr
	return nil
}

func op(y bool, ym map[string]any, f []string) error {
	if y {
		return operation_1(f, ym)
	}
	return operation_2(f, ym)
}

func operation(ym map[string]any, f []string) error {
	switch f[1] {
	case "==":
		return op(f[0] == f[2], ym, f)
	case "!=":
		return op(f[0] != f[2], ym, f)
	case "<":
		return op(unsafeConvert.StringToInt64(f[0]) < unsafeConvert.StringToInt64(f[2]), ym, f)
	case "<=":
		return op(unsafeConvert.StringToInt64(f[0]) <= unsafeConvert.StringToInt64(f[2]), ym, f)
	case ">":
		return op(unsafeConvert.StringToInt64(f[0]) > unsafeConvert.StringToInt64(f[2]), ym, f)
	case ">=":
		return op(unsafeConvert.StringToInt64(f[0]) >= unsafeConvert.StringToInt64(f[2]), ym, f)
	default:
		return Error_Invalid
	}
}

func JsonUrl(r []string) error {
	switch r[0] {
	case "parse":
		if len(r) != 4 {
			return Error_Invalid
		}
		vars[r[3]] = gjson.Get(JsonData[r[1]], r[2]).String()
	default:
		if len(r) != 2 {
			return Error_Invalid
		}
		if _, err := url.Parse(r[0]); err != nil {
			return err
		}

		resp := request(r[0])
		defer resp.RawBody().Close()

		if resp.StatusCode() != 200 {
			return Errorf("server returned status code: %v", (resp.StatusCode()))
		}

		rd := unsafeConvert.String(resp.Body())

		if rd != "" {
			JsonData[r[1]] = rd
		}
	}
	return nil
}

func mkdir(path string, mode ...fs.FileMode) error {
	return fsutil.Mkdir(path, mode...)
}

func remove(path string) error {
	return fsutil.Remove(path)
}

func touch(path string) error {
	f, err := fsutil.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func downloadFile(filepath string, url string) error {
	resp := request(url)
	defer resp.RawBody().Close()
	if resp.StatusCode() != 200 {
		return Errorf("connection failed! Server returned status code: %v\nUrl: %v\nUser-Agent: %v\n", resp.StatusCode(), resp.Request.URL, resp.RawResponse.Request.UserAgent())
	}

	Printf("GMake2: Connection info: %v\n", resp.Status())

	filename, _ := guessFilename(resp.RawResponse)
	if filepath != "." {
		filename = filepath
	}
	file, _ := fsutil.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	io.Copy(file, resp.RawBody())
	Printf("GMake2: Download saved to ./%v \n", filename)
	return nil
}

// Deprecated
//
// This method is about to be deprecated and is no longer supported
func (r *Req) Do(str ...string) error {
	switch str[0] {
	case "def":
		d := replace(str[2:])
		switch str[1] {
		case "header":
			r.Header, _ = json.TUnmarshalString[map[string]string](d)
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
			return Errors("Unknown method")
		}
	default:
		return r.Request()
	}
	// v := strings.ReplaceAll(strings.Trim(fmt.Sprint(str), "[]"), " ", " ")
	return nil
}

// Deprecated
//
// This method is about to be deprecated and is no longer supported
func (r *Req) Request() (err error) {
	if _, err := url.Parse(r.Uri); err != nil {
		return err
	}

	client := resty.NewWithClient(Client)

	if r.Header == nil {
		r.Header = make(map[string]string)
		r.Header = map[string]string{
			"User-Agent": UserAgent,
		}
	}
	if r.Header["User-Agent"] == "" {
		r.Header["User-Agent"] = UserAgent
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

	if err != nil {
		return err
	}

	defer r.Resp.RawBody().Close()

	if r.Resp.StatusCode() != 200 {
		return Errorf("server returned error code: %v", r.Resp.StatusCode())
	}
	Println("GMake2: @req: 200 ok")

	body := unsafeConvert.String(r.Resp.Body())
	if body != "" {
		if r.Value != "" {
			vars[r.Value] = body
		}
		Println(body)
	}
	return nil
}

func wait(v ...string) error {
	ar := len(v)
	if ar < 2 {
		return Errors("@wait bad format")
	}
	// v[0] print
	Println(replace(v[:ar-1]))
	var t string
	fmt.Print("=> ")
	if _, err := fmt.Scan(&t); err != nil {
		return err
	}
	vars[v[ar-1]] = t
	return nil
}
