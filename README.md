# GMake2
A make-like program, forked from https://github.com/fdxxw/gmake .

This branch extends some functionality.

[![GitHub Actions](https://github.com/3JoB/gmake2/actions/workflows/codeql.yml/badge.svg)](https://github.com/3JoB/gmake2/actions)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2F3JoB%2Fgmake2.svg?type=smail)](https://app.fossa.com/projects/git%2Bgithub.com%2F3JoB%2Fgmake2?ref=badge_smail)

# Menu

- [GMake2](#gmake2)
- [Menu](#menu)
- [Installing](#installing)
- [Getting Started](#getting-started)
- [Features](#features)
  - [Keywords](#keywords)
  - [Built-in variables](#built-in-variables)
    - [time.now](#timenow)
    - [time.utc](#timeutc)
    - [time.unix](#timeunix)
    - [time.utc\_unix](#timeutc_unix)
    - [runtime.os](#runtimeos)
    - [runtime.arch](#runtimearch)
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


## Built-in variables
For the convenience of use, gmake2 has some built-in available variables, which will be continuously updated.

### time.now
Current time

```
@echo {{.time.now}}
```

### time.utc
Current UTC time

```
@echo {{.time.utc}}
```

### time.unix
Current Unix Time

```
@echo {{.time.unix}}
```

### time.utc_unix
Current UTC Unix time

```
@echo {{.time.utc_unix}}
```

### runtime.os
Current system name

```
@echo {{.runtime.os}}
```

### runtime.arch
Current System Architecture

```
@echo {{.runtime.arch}}
```

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