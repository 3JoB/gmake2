package main

import (
	"fmt"
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
		// fmt.Println(pair[0])
		vars[pair[0]] = pair[1]
	}
	lines := strings.Split(cast.ToString(ym[commands]), "\n")
	for _, line := range lines {
		if line != "" {
			if strings.TrimSpace(line)[0] == '#' {
				continue
			}
			// line = ResolveVars(vars, line)
			cmdStrs, err := shellquote.Split(line)
			checkError(err)
			for i, cmdStr := range cmdStrs {
				cmdStrs[i] = ResolveVars(vars, cmdStr)
			}
			bin, args := cmdStrs[0], cmdStrs[1:]
			if len(args) == 0 {
				ErrPrintf("GMake2: Illegal instruction!\nGMake2: Error Command: %v \n", fmt.Sprint(cmdStrs[:]))
			}
			var BinMap = map[string]func(ym map[string]any, args []string) error{
				"@var":   KW_Var,
				"@env":   KW_Env,
				"@cmd":   KW_Run,
				"@echo":  KW_Echo,
				"@wait":  KW_Wait,
				"@end":   KW_End,
				"@if":    KW_Operation,
				"@sleep": KW_Sleep,
				"@val":   KW_Val,
				"#":      KW_Note,
				"@cd":    KW_Cd,
				"@touch": KW_Touch,
				"@mkdir": KW_Mkdir,
				"@mv":    KW_Mv,
				"@cp":    KW_Copy,
				"@rm":    KW_Del,
				"@req":   KW_Req,
				"@json":  KW_Json,
				"@dl":    KW_Downloads,
			}
			if fc, ok := BinMap[bin]; ok {
				checkError(fc(ym, args))
			} else {
				checkError(KW_Default(bin, args))
			}
		}
	}
}
