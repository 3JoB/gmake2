package main

var (
	VersionInfo string = `GMake2 is distributed under Mozilla Public License 2.0.
Github: https://github.com/3JoB/gmake2

Version: ` + SoftVersion + `
CommitID: ` + SoftCommit

	InitFileContent string = `config:
  default: all
  proxy:
  req: false

var:
  msg: GMake2
  
all: |
  @echo Hello! {{.msg}}!`
)
