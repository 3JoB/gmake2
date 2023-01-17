package main

var (
	VersionInfo string = `GMake2 is distributed under Apache-2.0 license.
Github: https://github.com/3JoB/gmake2

Version: ` + SoftVersion + `
CommitID: ` + SoftCommit

	InitFileContent string = `config:
  default: all

var:
  msg: GMake2
  
all: |
  @echo Hello! {{.msg}}!`
)