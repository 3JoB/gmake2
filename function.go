package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func KW_Note(ym map[string]any, args []string) error {
	return nil
}

func KW_Var(ym map[string]any, args []string) error {
	vars[args[0]] = strings.Join(args[1:], " ")
	return nil
}

func KW_Env(ym map[string]any, args []string) error {
	os.Setenv(args[0], strings.Join(args[1:], " "))
	return nil
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
	time.Sleep(time.Second * cast.ToDuration(args[0]))
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
	abs, err := filepath.Abs(args[0])
	cmdDir = abs
	return err
}

func KW_Mv(ym map[string]any, args []string) error {
	copy(args[0], args[1])
	os.RemoveAll(args[0])
	return nil
}

func KW_Copy(ym map[string]any, args []string) error {
	copy(args[0], args[1])
	return nil
}

func KW_Del(ym map[string]any, args []string) error {
	return os.RemoveAll(args[0])
}

func KW_Mkdir(ym map[string]any, args []string) error {
	mkdir(args[0])
	return nil
}

func KW_Touch(ym map[string]any, args []string) error {
	touch(args[0])
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
		R.Network(args...)
	} else {
		return errors.New("GMake2: The @req tag has been deprecated.")
	}
	return nil
}

func KW_Default(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	if cmdDir != "" {
		cmd.Dir = cmdDir
	}
	ExecCmd(cmd)
	return nil
}