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
    - [@env](#env)
    - [@cmd](#cmd)
    - [comment](#comment)
    - [@touch](#touch)
    - [@mv](#mv)
    - [@copy](#copy)
    - [@rm](#rm)
    - [@mkdir](#mkdir)
    - [@cd](#cd)
  - [system command](#system-command)
- [examples](#examples)

# Installing

Download the latest version from github.


```sh
go get -u github.com/3JoB/gmake2
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

## system command

System commands, execute console commands, and execute everything that the console can execute.

```sh
go build
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
