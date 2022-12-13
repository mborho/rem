# rem - remember

[![Go](https://github.com/mborho/rem/actions/workflows/go.yml/badge.svg)](https://github.com/mborho/rem/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/mborho/rem/badge.svg?branch=master)](https://coveralls.io/github/mborho/rem?branch=master)
[![rem](https://snapcraft.io/rem/badge.svg)](https://snapcraft.io/rem)

rem is a small tool for remembering things on the command line. It can be used to store commands and execute them later on. Or to simply store lines of text.

### Usage
```sh
$ rem [flags] command [argument]
```

### Commands

*    -h, help - Shows this help.
*    -a, add [string] - Adds a command/text.
*    rm [index] - Removes line with given index number.
*    echo [index] - Displays line with given index number.
*    -f, filter [regexp] - Filters stored commands by given regular expression.
*    here - Creates a .rem file in the given directory. Default: **~/.rem**
*    clear - Clears currently active .rem file, **./.rem** or **~/.rem**
*    [index|tag] - Executes line with given index number / tag name.

### Flags

* -g - Use global rem file ~/.rem
* -t - Tag for command when adding with -a/add.
* -p - Print command to stdout before executing index/tag.


Run **rem** without any arguments to list all stored commands/strings.

### Install
```
$ go get github.com/mborho/rem
```

### Install binary 

On Linux you can install a released binary directly:

```
$ sudo wget https://github.com/mborho/rem/releases/download/v0.16.2/rem_0.16.2_linux_amd64 \
    -O /usr/local/bin/rem && sudo chmod +x /usr/local/bin/rem
```
 See [releases](https://github.com/mborho/rem/releases) for the specific binary (amd64, arm, arm64) to use.
 
### Snap package

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/rem)

### Examples

Create a local .rem file in the directory you are in:
```sh
$ rem here
```
Remember commands:
```sh
$ rem add ls -la
$ rem add mysqldump -u name -h 172.17.42.1 -P 49176 -p demo-db
```
List commands stored:
```sh
$ rem
 0  ls -la
 1  mysqldump -u name -h 172.17.42.1 -P 49176 -p demo-db
```
Execute **ls -la** (listnumber 0)
```sh
$ rem 0
insgesamt 12
drwxrwxr-x 2 martin martin 4096 Dez  5 14:41 .
drwxrwxr-x 6 martin martin 4096 Dez  5 14:33 ..
-rw------- 1 martin martin   60 Dez  5 14:42 .rem
```

You can also use parenthesis like **$(pwd)** or operators like **&&** or pipes **|**. If neccessary usse quotes to escape when adding a command.

```sh
$ rem -t count-lines -a 'cat $HOME/.bashrc | wc'
$ rem count-lines
    179     709    6053
```

Remove a command:
```sh
$ rem rm 1
```

### Demo

[![asciicast](https://asciinema.org/a/pvaQM8E5CGYJTPSQ4RhiWotEi.svg)](https://asciinema.org/a/pvaQM8E5CGYJTPSQ4RhiWotEi)

### Thanks

[psolbach](https://github.com/psolbach) for the name **rem**!

### License

GPLv3 - see [LICENSE](https://raw.githubusercontent.com/mborho/rem/master/LICENSE)

