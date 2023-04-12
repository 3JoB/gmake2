package main

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/3JoB/ulib/fsutil"
	"github.com/3JoB/ulib/path"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

/*func Operation[T any](e bool, trueValue, falseValue T) T {
	if e {
		return trueValue
	}
	return falseValue
}*/

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
	return fsutil.TruncWrite(path, v)
}

func replace(v []string) string {
	return strings.ReplaceAll(strings.Trim(fmt.Sprint(v), "[]"), " ", " ")
}

func split(v, r string) []string {
	return strings.Split(v, r)
}

func checkError(err error) {
	if err != nil {
		ErrPrintf("GMake2: %v", err.Error())
	}
}

func Println(a ...any) {
	fmt.Println(a...)
}

func Printf(f string, v ...any) {
	fmt.Printf(f, v...)
}

func Sprintf(f string, v ...any) string {
	return fmt.Sprintf(f, v...)
}

func Sprintln(v ...any) string {
	return fmt.Sprintln(v...)
}

func Errorf(format string, v ...any) error {
	return fmt.Errorf(format, v...)
}

func ErrPrintln(a ...any) {
	if debug {
		panic(Sprintln(a...))
	}
	Println(a...)
	Exit()
}

func ErrPrintf(format string, v ...any) {
	if debug {
		panic(Sprintf(format, v...))
	}
	Printf(format, v...)
	Exit()
}

func E(c BinConfig, err error) {
	if err != nil {
		if c.YamlData == nil {
			ErrPrintln(ErrCommand(c.CommandLine, c.CommandGroup, err, ""))
		} else {
			ErrPrintln(ErrCommand(c.CommandLine, c.CommandGroup, err, c.YamlDataBin+" "+replace(c.YamlData)))
		}
	}
}

func Errors(errs string) error {
	return errors.New(errs)
}

func ErrCommand(line int, group, msg any, command string) error {
	line++
	if command == "" {
		return Errorf("gmake2: %v\n    Errored command group: %v\n    Errored row count: %v", msg, group, line)
	}
	return Errorf("gmake2: %v\n    Errored command group: %v\n    Errored command: %v\n    Errored row count: %v", msg, group, command, line)
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

	filename = path.Base(path.Clean("/" + filename))
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
