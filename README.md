# GMake2

<p align="center">
    <p align="center"><img src="wiki/gmake2.png"></p>
    <p align="center">This image is from <a href="https://quasilyte.dev/gopherkon/">Gopher Konstructor</a></p>
    <p align="center"><strong>Build a GMakefile at lightning speed!</strong></p>
    <p align="center">
        <a href="https://github.com/3JoB/gmake2/actions"><img src="https://img.shields.io/github/actions/workflow/status/3JoB/gmake2/codeql.yml?label=CodeQL%20Scanner&style=flat-square" alt="GitHub Workflow Status"></a>
        <a href="https://app.fossa.com/projects/git%2Bgithub.com%2F3JoB%2Fgmake2?ref=badge_smail"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2F3JoB%2Fgmake2.svg?type=smail" alt="FOSSA Status"></a>
        <a href="https://github.com/3JoB/gmake2/blob/master/LICENSE"><img src="https://img.shields.io/github/license/3JoB/gmake2?style=flat-square" alt="MPL-2.0"></a>
        <a href="#"><img src="https://img.shields.io/github/go-mod/go-version/3JoB/gmake2?label=Go%20Version&style=flat-square" alt="Go Version"></a>
        <a href="https://github.com/3JoB/gmake2/release"><img src="https://img.shields.io/github/v/release/3JoB/gmake2?label=Release%20Version&style=flat-square" alt="GitHub release (latest by date)"></a>
    </p>
    <p align="center">
        <a href="https://github.com/3JoB/gmake2/issues"><img src="https://img.shields.io/github/issues/3JoB/gmake2?label=GMake2%20Issues&style=flat-square" alt="GitHub Issues"></a>
        <a href="https://github.com/3JoB/gmake2/stargazers"><img src="https://img.shields.io/github/stars/3JoB/gmake2?label=Stars&style=flat-square" alt="GitHub Repo stars"></a>
        <a href="#"><img src="https://img.shields.io/github/downloads/3JoB/gmake2/latest/total?label=Downloads%40Latest&style=flat-square" alt="GitHub release (latest by date)"></a>
        <a href="#"><img src="https://img.shields.io/github/repo-size/3JoB/gmake2?style=flat-square" alt="GitHub repo size"></a>
        <a href="#"><img src="https://img.shields.io/github/commit-activity/m/3JoB/gmake2?style=flat-square" alt="GitHub commit activity"></a>
    </p>
</p>


This project is the follow-up maintenance of [go-gmake](https://github.com/fdxxw/gmake).


# Menu

- [GMake2](#gmake2)
- [Menu](#menu)
- [Installing](#installing)
  - [Install from software source](#install-from-software-source)
  - [Install using Chocolatey \[Windows\]](#install-using-chocolatey-windows)
  - [Install from Github Releases](#install-from-github-releases)
  - [Install from source code](#install-from-source-code)
- [Getting Started](#getting-started)
- [Features](#features)
  - [Keywords](#keywords)
  - [Reserved Keyword](#reserved-keyword)
  - [Variables](#variables)
- [Examples](#examples)
- [Other information](#other-information)
- [License](#license)

# Installing

## Install from software source
This method is limited to systems using apt and dpkg for package management.


Please execute the following commands in the order they were written.
```sh
echo 'deb https://deb.lcag.org stable main' | sudo tee /etc/apt/sources.list.d/malonan.list

wget -qO - https://deb.lcag.org/public.key | sudo apt-key --keyring /etc/apt/trusted.gpg.d/malonan.gpg add -

sudo apt update && sudo apt install gmake2
```

Upgrade GMake2
```sh
apt update && apt upgrade
```


<strong>The following method is proven to be no longer available, it causes apt not to read the installed key, the workaround is to delete the key and re-add it using apt-key.</strong>

```sh
wget -qO - https://deb.lcag.org/public.key | sudo gpg --no-default-keyring --keyring gnupg-ring:/etc/apt/trusted.gpg.d/malonan.gpg --import
```

## Install using Chocolatey [Windows]

chocolatey is a package manager on windows, windows users can quickly install GMake2 through `choco` command.

Since the review of chocolatey takes a long time, it takes about a day for chocolatey to be updated synchronously after each update of GMake2.

Chocolatey only provides GMake2 for amd64 architecture.

```sh
choco install gmake2 --version=[latest version]

# Example: GMake2 v2.3.0-LTS
choco install gmake2 --version=2.3.0

```

## Install from Github Releases
Download the latest version from github.


[Release](https://github.com/3JoB/gmake2/releases)

## Install from source code
You can build gmake2 directly using Go build, but the version subcommand will not work properly.

```sh
git clone https://github.com/3JoB/gmake2 && cd gmake2

# gmake2 installed
gmake2

# gmake2 is not installed
export CGO_ENABLED=0
go build -ldflags "-s -w -X 'main.SoftCommit=owner' -X 'main.SoftVersion=owner'"
```



# Getting Started

Write the GMakefile.yml file in the current directory, the content is as follows

```yml
vars:
  msg: Hello World

all: |
  @echo {{.msg}}

mg: |
  @echo What's up???
```

Then run `gmake2` on the current command line console, you can see the console print

```
Hello World
```

Or execute `gmake2 mg` to execute the specified command, and the console will print
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


# Examples

GMakefile.yml

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
  @cp to.txt from.txt
  @rm from.txt
  @rm to.txt
  @mkdir from
  @mv from to
  @cp to from
  @rm from
  @rm to
  @env GOOS linux
  go build
```

```sh
gmake2
```

# Other information
- Due to unsigned and some high-risk operations, GMake2 may be blocked by some anti-virus software. If you have installed anti-virus software on your device, please manually set GMake2 to the whitelist.
- The binary released by GMake2 is compiled directly from the git library, the specific steps: 
  - 1_ Write the update 
  - 2_ Push to github 
  - 3_ Run the gmake2 command to compile the binary
- GMake2 only has the following binary distribution channels, and cannot guarantee that other channels are safe: 
  - [DEB Server](https://deb.lcag.org)
  - [Chocolatey](https://lcag.org/gmake2.choco) 
  - [Github Release](https://lcag.org/gmake2.releases). 
- If GMake2 is installed using apt or chocolatey, GMake2's built-in self-update command cannot be used unless the `--upgrade` tag is used to force an update

# License
This software is distributed under MPL-2.0 license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2F3JoB%2Fgmake2.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2F3JoB%2Fgmake2?ref=badge_large)
