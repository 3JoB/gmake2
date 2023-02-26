package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

type HandlerFunc func(ym map[string]any, args []string) error

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
func KW_End(ym map[string]any, args []string) error {
	Exit()
	return nil
}

// Note
func KW_Note(ym map[string]any, args []string) error {
	return nil
}

// Add/override GMakefile variables
func KW_Var(ym map[string]any, args []string) error {
	vars[args[0]] = strings.Join(args[1:], " ")
	return nil
}

// Set terminal environment variables
func KW_Env(ym map[string]any, args []string) error {
	if len(args) < 2 {
		return exec.ErrNotFound
	}
	return os.Setenv(args[0], strings.Join(args[1:], " "))
}

func KW_Run(ym map[string]any, args []string) error {
	run(ym, args[0])
	return nil
}

func KW_Wait(ym map[string]any, args []string) error {
	wait(args...)
	return nil
}

func KW_Sleep(ym map[string]any, args []string) error {
	if len(args) != 1 {
		return ErrCommand()
	}
	sleep(args[0])
	return nil
}

func KW_Operation(ym map[string]any, args []string) error {
	return operation(ym, args)
}

func KW_Val(ym map[string]any, args []string) error {
	arg := args[2:]
	cmd := exec.Command(args[1], arg...)
	if cmdDir != "" {
		cmd.Dir = cmdDir
	}
	val(args, cmd)
	return nil
}

func KW_Echo(ym map[string]any, args []string) error {
	Println(strings.Join(args, " "))
	return nil
}

func KW_Cd(ym map[string]any, args []string) error {
	abs, err := filepath.Abs(replace(args))
	cmdDir = abs
	return err
}

func KW_Mv(ym map[string]any, args []string) error {
	copy(args[0], args[1])
	remove(args[0])
	return nil
}

func KW_Copy(ym map[string]any, args []string) error {
	copy(args[0], args[1])
	return nil
}

func KW_Del(ym map[string]any, args []string) error {
	remove(replace(args))
	return nil
}

func KW_Mkdir(ym map[string]any, args []string) error {
	mkdir(replace(args))
	return nil
}

func KW_Touch(ym map[string]any, args []string) error {
	touch(replace(args))
	return nil
}

func KW_Json(ym map[string]any, args []string) error {
	return JsonUrl(args)
}

func KW_Downloads(ym map[string]any, args []string) error {
	if len(args) == 1 {
		downloadFile(".", args[0])
	} else {
		downloadFile(args[1], args[0])
	}
	return nil
}

func KW_Req(ym map[string]any, args []string) error {
	if cast.ToBool(cfg["req"]) {
		R.Do(args...)
	} else {
		ErrPrintln("GMake2: The @req tag has been deprecated.")
	}
	return nil
}

func KW_Hash(ym map[string]any, args []string) error {
	return nil
}

func KW_Default(bin string, args []string) error {
	switch bin[0:1] {
	case "#":
	case "@":
		ErrPrintf("GMake2: GMake2 keyword %v unregistered", bin)
	default:
		cmd := exec.Command(bin, args...)
		if cmdDir != "" {
			cmd.Dir = cmdDir
		}
		ExecCmd(cmd)
	}
	return nil
}
