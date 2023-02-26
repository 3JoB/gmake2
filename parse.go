package main

import (
	"bytes"
	"os"
	"runtime"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

func variable(v map[string]any) map[string]any {
	v["time"] = map[string]any{
		"now":      time.Now().Format("2006-01-02 15:04"),
		"utc":      time.Now().UTC().Format("2006-01-02 15:04"),
		"unix":     time.Now().Unix(),
		"utc_unix": time.Now().UTC().Unix(),
	}
	v["runtime"] = map[string]any{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
	v["gmake2"] = map[string]any{
		"version": SoftVersion,
		"code":    SoftVersionCode,
		"time":    SoftBuildTime,
	}
	return v
}

func parseConfig(cfgFile string) map[string]any {
	ymlData, err := os.ReadFile(cfgFile)
	checkError(err)
	m := make(map[string]any)
	err = yaml.Unmarshal(ymlData, &m)
	checkError(err)
	return m
}

func parseMap(ym map[string]any) {
	if v, ok := ym["vars"]; ok {
		vars = v.(map[string]any)
	} else {
		vars = make(map[string]any)
	}
	if v, ok := ym["config"]; ok {
		cfg = v.(map[string]any)
	} else {
		cfg = make(map[string]any)
	}
	vars = variable(vars)
}

func parseTags(v string) {
	tags := split(v, ",")
	if len(tags) == 0 {
		return
	}
	for _, tag := range tags {
		data := split(tag, "=")
		if len(data) == 0 {
			continue
		}
		vars[data[0]] = data[1]
	}
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

/*func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(`GMake2: Unclosed quote in command line: " ` + command + ` "`)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
*/
