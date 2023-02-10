# GMake2 KeyWord

Keywords are some built-in instructions of GMake2, which can be used to quickly build GMakefile.

# Menu

- [GMake2 KeyWord](#gmake2-keyword)
- [Menu](#menu)
    - [comment](#comment)
    - [@cmd](#cmd)
    - [@cd](#cd)
    - [@download](#download)
    - [@echo](#echo)
    - [@env](#env)
    - [@wait](#wait)
    - [@if](#if)
    - [@json](#json)
    - [@var](#var)
    - [@val](#val)
    - [@sleep](#sleep)
    - [@touch](#touch)
    - [@mv](#mv)
    - [@copy](#copy)
    - [@rm](#rm)
    - [@mkdir](#mkdir)
  - [system command](#system-command)


### comment

`#` Begins with comments

```
# comment
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


### @echo

Print information

```sh
@echo msg
```

### @env

Set environment variables

```
@env GOOS linux
```

### @wait
Wait for the command line to enter the value, and then assign the value.

Example: `@wait [prompt information] [value]`

next demo:

GMakefile.yml
```yml
hello: |
  @wait Please Enter Password pass
  @echo {{.pass}}
  @echo 114
```

CommandLine:
```sh
$ gmake2 hello

[Please Enter Password]
=> 

;Enter 2023;

[Please Enter Password]
=> 2023
2023
114
```


### @if
This keyword is a binary operator.

Due to the long length, please go to [if.md](if.md) to view.

### @json
Obtain json information from a link, and assign values to variables after obtaining keywords.

example.json
```json
{
    "name": "make!"
}
```

GMakefile
```yml
all: |
  @json url https://example.com/example.json string name myname
  @echo {{.myname}}
```

The following are the data types supported by `@json` keyword:
```
string  (or String)
bool    (or Bool)
int64   (or int,int8,in16,in32)
uint64  (or uint,uint8.uint16,uint32)
float64 (or float,float32)
```

When the data type is set incorrectly, the `String` type will be automatically used.


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

### @val
Execute a custom command and assign the returned value to a variable.

```yml
all: |
  @val commit git log --pretty=format:'%h' -1
  @echo {{.commit}}
```


### @sleep
Let Gmake2 be paused to run.

Example: Pause five seconds
```yml
@sleep 5
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


## system command

System commands, execute console commands, and execute everything that the console can execute.

```sh
go build
```