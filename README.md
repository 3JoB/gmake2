# GMake2
A make-like program, forked from https://github.com/fdxxw/gmake .

This branch extends some functionality.

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/3JoB/gmake2/codeql.yml?label=CodeQL%20Scanner&style=flat-square)](https://github.com/3JoB/gmake2/actions)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2F3JoB%2Fgmake2.svg?type=smail)](https://app.fossa.com/projects/git%2Bgithub.com%2F3JoB%2Fgmake2?ref=badge_smail)
[![Apache-2.0](https://img.shields.io/github/license/3JoB/gmake2?style=flat-square)](https://github.com/3JoB/gmake2/blob/master/LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/3JoB/gmake2?label=Go%20Version&style=flat-square)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/3JoB/gmake2?label=Release%20Version&style=flat-square)](https://github.com/3JoB/gmake2/release)

[![GitHub issues](https://img.shields.io/github/issues/3JoB/gmake2?label=GMake2%20Issues&style=flat-square)](https://github.com/3JoB/gmake2/issues)
[![GitHub Repo stars](https://img.shields.io/github/stars/3JoB/gmake2?label=Stars&style=flat-square)](https://github.com/3JoB/gmake2/stargazers)
![GitHub release (latest by date)](https://img.shields.io/github/downloads/3JoB/gmake2/latest/total?label=Downloads%40Latest&style=flat-square)
![GitHub repo size](https://img.shields.io/github/repo-size/3JoB/gmake2?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/3JoB/gmake2?style=flat-square)


# Menu

- [GMake2](#gmake2)
- [Menu](#menu)
- [Installing](#installing)
- [Getting Started](#getting-started)
- [Features](#features)
  - [Keywords](#keywords)
  - [Reserved Keyword](#reserved-keyword)
  - [Variables](#variables)
- [examples](#examples)
- [License](#license)

# Installing

Download the latest version from github.


```sh
https://github.com/3JoB/gmake2/releases
```

# Getting Started

Write the gmake.yml file in the current directory, the content is as follows

```yml
vars:
  msg: Hello World

all: |
  @echo {{.msg}}

mg: |
  @echo What's up???
```

Then run `gmake` on the current command line console, you can see the console print

```
Hello World
```

Or execute `gmake mg` to execute the specified command, and the console will print
```
What's up???
```
<font color=#e40d0d size=5>gmake2 automatically selects the all command when no command is specified.</font>
<br>

# Features

## Keywords

Keywords moved to [Wiki](wiki/Keyword.md)

## Reserved Keyword
View in [Wiki](wiki/Reserved_Keyword.md)


## Variables
Variables moved to [Wiki](wiki/variables.md)


# examples

examples.yml

```yml
vars:
  msg: Hello World
all: |
  @echo {{.msg}}
  # Modify the msg variable
  @var msg Hello
  @echo {{.msg}}
  # Create a file
  @touch from.txt
  @mv from.txt to.txt
  @copy to.txt from.txt
  @rm from.txt
  @rm to.txt
  @mkdir from
  @mv from to
  @copy to from
  @rm from
  @rm to
  @env GOOS linux
  go build
```

```sh
gmake
```

# License
This software is distributed under Apache-2.0 license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2F3JoB%2Fgmake2.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2F3JoB%2Fgmake2?ref=badge_large)