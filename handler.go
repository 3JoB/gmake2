package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/3JoB/ulib/fsutil"
	"github.com/3JoB/unsafeConvert"
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
		c.Error(Error_Invalid)
	}
	c.Error(os.Setenv(c.YamlData[0], strings.Join(c.YamlData[1:], " ")))
}

func KW_Run(c BinConfig) {
	if c.YamlDataLine != 1 {
		c.Error(Error_Invalid)
	}
	run(c.YamlConfig, c.YamlData[0])
}

func KW_Wait(c BinConfig) {
	c.Error(wait(c.YamlData...))
}

func KW_Sleep(c BinConfig) {
	if c.YamlDataLine != 1 {
		c.Error(Error_Invalid)
	}
	time.Sleep(time.Second * time.Duration(unsafeConvert.StringToInt64(c.YamlData[0])))
}

func KW_Operation(c BinConfig) {
	c.Error(operation(c.YamlConfig, c.YamlData))
}

func KW_Val(c BinConfig) {
	arg := c.YamlData[2:]
	cmd := exec.Command(c.YamlData[1], arg...)
	if cmdDir != "" {
		cmd.Dir = cmdDir
	}
	c.Error(val(c.YamlData, cmd))
}

func KW_Echo(c BinConfig) {
	Println(strings.Join(c.YamlData, " "))
}

func KW_Cd(c BinConfig) {
	if abs, err := filepath.Abs(replace(c.YamlData)); err != nil {
		c.Error(err)
	} else {
		cmdDir = abs
	}
}

func KW_Mv(c BinConfig) {
	c.Error(fsutil.CopyAll(c.YamlData[0], c.YamlData[1]))
	c.Error(remove(c.YamlData[0]))
}

func KW_Copy(c BinConfig) {
	c.Error(fsutil.CopyAll(c.YamlData[0], c.YamlData[1]))
}

func KW_Del(c BinConfig) {
	c.Error(remove(replace(c.YamlData)))
}

func KW_Mkdir(c BinConfig) {
	c.Error(mkdir(replace(c.YamlData)))
}

func KW_Touch(c BinConfig) {
	c.Error(touch(replace(c.YamlData)))
}

func KW_Json(c BinConfig) {
	c.Error(JsonUrl(c.YamlData))
}

func KW_Downloads(c BinConfig) {
	if c.YamlDataLine == 1 {
		c.Error(downloadFile(".", c.YamlData[0]))
	} else {
		c.Error(downloadFile(c.YamlData[1], c.YamlData[0]))
	}
}

func KW_Req(c BinConfig) {
	if cast.ToBool(cfg["req"]) {
		c.Error(R.Do(c.YamlData...))
	} else {
		c.Error(Errors("GMake2: The @req tag has been deprecated."))
	}
}

func KW_Hash(c BinConfig) {}

func KW_Default(bin string, c BinConfig) {
	switch bin[0:1] {
	case "#":
	case "@":
		c.Error(Errors(Sprintf("GMake2 keyword %v unregistered", bin)))
	default:
		cmd := exec.Command(bin, c.YamlData...)
		if cmdDir != "" {
			cmd.Dir = cmdDir
		}
		c.Error(ExecCmd(cmd))
	}
}
