package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/3JoB/ulib/fsutil"
	"github.com/spf13/cast"
)

type HandlerFunc func(c BinConfig)

type BinConfig struct {
	YamlConfig   map[string]any // Raw data of command group
	YamlData     []string       // Commands within a command group
	YamlDataBin  string         // Yaml Data Bin
	YamlDataLine int            // The number of lines of the command in the command group
	CommandGroup string         // The name of the group where the command is located
	CommandLine  int            // The line where the command is
}

var BinMap map[string]HandlerFunc

func init() {
	BinMap = map[string]HandlerFunc{
		"@var":   KW_Var,
		"@env":   KW_Env,
		"@run":   KW_Run,
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
		"@hash":  KW_Hash,
	}
}

// Terminate operation
func KW_End(c BinConfig) {
	Exit()
}

// Note
func KW_Note(c BinConfig) {}

// Add/override GMakefile variables
func KW_Var(c BinConfig) {
	vars[c.YamlData[0]] = strings.Join(c.YamlData[1:], " ")
}

// Set terminal environment variables
func KW_Env(c BinConfig) {
	if c.YamlDataLine < 2 {
		E(c, Error_Invalid)
	}
	E(c, os.Setenv(c.YamlData[0], strings.Join(c.YamlData[1:], " ")))
}

func KW_Run(c BinConfig) {
	if c.YamlDataLine != 1 {
		E(c, Error_Invalid)
	}
	run(c.YamlConfig, c.YamlData[0])
}

func KW_Wait(c BinConfig) {
	E(c, wait(c.YamlData...))
}

func KW_Sleep(c BinConfig) {
	if c.YamlDataLine != 1 {
		E(c, Error_Invalid)
	}
	sleep(c.YamlData[0])
}

func KW_Operation(c BinConfig) {
	E(c, operation(c.YamlConfig, c.YamlData))
}

func KW_Val(c BinConfig) {
	arg := c.YamlData[2:]
	cmd := exec.Command(c.YamlData[1], arg...)
	if cmdDir != "" {
		cmd.Dir = cmdDir
	}
	E(c, val(c.YamlData, cmd))
}

func KW_Echo(c BinConfig) {
	Println(strings.Join(c.YamlData, " "))
}

func KW_Cd(c BinConfig) {
	if abs, err := filepath.Abs(replace(c.YamlData)); err != nil {
		E(c, err)
	} else {
		cmdDir = abs
	}
}

func KW_Mv(c BinConfig) {
	E(c, fsutil.CopyAll(c.YamlData[0], c.YamlData[1]))
	E(c, remove(c.YamlData[0]))
}

func KW_Copy(c BinConfig) {
	E(c, fsutil.CopyAll(c.YamlData[0], c.YamlData[1]))
}

func KW_Del(c BinConfig) {
	E(c, remove(replace(c.YamlData)))
}

func KW_Mkdir(c BinConfig) {
	E(c, mkdir(replace(c.YamlData)))
}

func KW_Touch(c BinConfig) {
	E(c, touch(replace(c.YamlData)))
}

func KW_Json(c BinConfig) {
	E(c, JsonUrl(c.YamlData))
}

func KW_Downloads(c BinConfig) {
	if c.YamlDataLine == 1 {
		E(c, downloadFile(".", c.YamlData[0]))
	} else {
		E(c, downloadFile(c.YamlData[1], c.YamlData[0]))
	}
}

func KW_Req(c BinConfig) {
	if cast.ToBool(cfg["req"]) {
		E(c, R.Do(c.YamlData...))
	} else {
		E(c, Errors("GMake2: The @req tag has been deprecated."))
	}
}

func KW_Hash(c BinConfig) {}

func KW_Default(bin string, c BinConfig) {
	switch bin[0:1] {
	case "#":
	case "@":
		E(c, Errors(Sprintf("GMake2 keyword %v unregistered", bin)))
	default:
		cmd := exec.Command(bin, c.YamlData...)
		if cmdDir != "" {
			cmd.Dir = cmdDir
		}
		E(c, ExecCmd(cmd))
	}
}
