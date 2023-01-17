package main

import (
	"bytes"
	"errors"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

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

func parseCommandLine(command string) ([]string, error) {
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
