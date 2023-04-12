package main

import (
	"net/http"
	"os"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cast"
)

var (
	vars     map[string]any
	cfg      map[string]any
	JsonData map[string]string
	R        Req
	Client   *http.Client
	debug    bool
	cmdDir   string
)

func init() {
	JsonData = make(map[string]string)
}

func run(ym map[string]any, commands string) {
	if ym[commands] == nil {
		ErrPrintf("GMake2: Command not found %v\n", commands)
	}

	if commands == "vars" || commands == "config" {
		ErrPrintf("GMake2: Illegal command group name!")
	}

	if commands == "end" {
		Exit()
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		vars[pair[0]] = pair[1]
	}

	// lines := strings.Split(cast.ToString(ym[commands]), "\n")
	lines := ym[commands].([]any)
	for is, line := range lines {
		if line != nil {
			if strings.TrimSpace(cast.ToString(line))[0] == '#' {
				continue
			}
			// line = ResolveVars(vars, line)
			cmdStrs, err := shellquote.Split(cast.ToString(line))
			BinConfig{CommandLine: is, CommandGroup: commands}.Error(err)
			for i, cmdStr := range cmdStrs {
				cmdStrs[i] = ResolveVars(vars, cmdStr)
			}
			bin, args := cmdStrs[0], cmdStrs[1:]
			if len(args) == 0 {
				ErrPrintf("GMake2: Illegal instruction!\nGMake2: Error Command: %v \n", strings.Join(cmdStrs, " "))
			}
			bins := BinConfig{
				YamlConfig:   ym,
				YamlData:     args,
				YamlDataLine: len(args),
				YamlDataBin:  bin,
				CommandGroup: commands,
				CommandLine:  is,
			}
			if fc, ok := BinMap[bin]; ok {
				fc(bins)
			} else {
				KW_Default(bin, bins)
			}
		}
	}
}
