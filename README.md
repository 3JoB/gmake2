# GMake2
A make-like program, forked from https://github.com/fdxxw/gmake .

This branch extends some functionality.

# Menu

- [GMake2](#gmake2)
- [Menu](#menu)
- [Installing](#installing)
- [Getting Started](#getting-started)
- [Features](#features)
  - [build-in command](#build-in-command)
    - [@echo](#echo)
    - [@var](#var)
    - [@if](#if)
    - [@env](#env)
    - [@cmd](#cmd)
    - [comment](#comment)
    - [@touch](#touch)
    - [@mv](#mv)
    - [@copy](#copy)
    - [@rm](#rm)
    - [@mkdir](#mkdir)
    - [@cd](#cd)
    - [@download](#download)
  - [system command](#system-command)
  - [Built-in variables](#built-in-variables)
    - [time](#time)
    - [time\_utc](#time_utc)
    - [time\_unix](#time_unix)
    - [time\_utc\_unix](#time_utc_unix)
    - [runtime\_os](#runtime_os)
    - [runtime\_arch](#runtime_arch)
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

## build-in command

The built-in command is as follows

### @echo

Print information

```sh
@echo msg
```

### @var

Set variable

```sh
@var msg Hello World
```

or

```yml
vars:
  msg: Hello World

all: |
  @echo {{.msg}}
```

### @if
This keyword is a binary operator that supports the following two data types:
`string, int64`

<strong>Example</strong>

```yml
# Equal to
all: |
  @if windows == windows then iaw

iaw: |
  @echo i am windows!!!


# or
all: |
  @if linux == windows then iaw or ial

iaw: |
  @echo i am windows!!!

ial: |
  @echo i am linux!!!


# Not equal to
all: |
  @if windows2 != windows then iaw

iaw: |
  @echo i am not windows!!!


# Greater than
all: |
  @if 2 > 1 then iaw

iaw: |
  @echo i am 2!!!

# Greater than or equal
all: |
  @if 2 >= 2 then iaw

iaw: |
  @echo i am 2!!!

# Smaller than
all: |
  @if 1 < 2 then iaw

iaw: |
  @echo i am 1!!!

# Less than or equal to
all: |
  @if 1 <= 1 then iaw

iaw: |
  @echo i am 1!!!
```

### @env

Set environment variables

```
@env GOOS linux
```

### @cmd
Execute other configured commands.

Example:
```yml
all: |
  @cmd readme

readme: |
  @echo I Like It
```

### comment

`#` Begins with comments

```
# comment
```

### @touch

Create a file

```
@touch from.txt
```

### @mv

Move a file or directory

```
@mv from.txt to.txt
```

### @copy

Copy a file or directory

```
@copy to.txt from.txt
```

### @rm

Delete file or directory

```
@rm from.txt
```

### @mkdir

Create a directory

```
@mkdir from
```

### @cd

Set the directory to make subsequent console commands run in the specified directory. It is only valid for system commands and invalid for built-in commands.

```
@cd from
```

### @download
Download a file from the server.

Example: save the file as-is to the current directory.
```sh
@download https://github.com/3JoB/gmake2/releases/download/v2.0.0/gmake2_Linux_aarch64.tar.gz
```

Or

Example: save a file into a custom directory/file.
```sh
@download https://github.com/3JoB/gmake2/releases/download/v2.0.0/gmake2_Linux_aarch64.tar.gz bin/gmake2.tar.gz
```

## system command

System commands, execute console commands, and execute everything that the console can execute.

```sh
go build
```

## Built-in variables
For the convenience of use, gmake2 has some built-in available variables, which will be continuously updated.

### time
Current time

```
@echo {{.time}}
```

### time_utc
Current UTC time

```
@echo {{.time}}
```

### time_unix
Current Unix Time

```
@echo {{.time_unix}}
```

### time_utc_unix
Current UTC Unix time

```
@echo {{.time_utc_unix}}
```

### runtime_os
Current system name

```
@echo {{.runtime_os}}
```

### runtime_arch
Current System Architecture

```
@echo {{.runtime_arch}}
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