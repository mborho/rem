# rem - remember

rem is a small tool for remembering things on the command line. It can be used to store commands and execute them later on. Or to simply store lines of text.

### Usage
```sh
$ rem command [argument]
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


Run **rem** without any arguments to list all stored commands/strings.

### Build & Run

```sh
$ git clone https://github.com/mborho/rem
$ cd rem
$ GOPATH=`pwd` go build -o rem main.go line.go
$ ./rem
```

### PPA ##
rem is installable on Ubuntu x86_64 via [ppa:martin-borho/rem](https://launchpad.net/~martin-borho/+archive/ubuntu/rem)


```
#!bash

sudo add-apt-repository ppa:martin-borho/rem
sudo apt-get update
sudo apt-get install rem
```

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
It's not possible to execute commands, who use parenthesis like **$(pwd)**, operators like **&&** or pipes **|**. Just simple commands. If **rem** exeuctes a command, the current **rem** process will be replaced by the command executed.

Remove a command:
```sh
$ rem rm 1
```

### Thanks

[psolbach](https://github.com/psolbach) for the name **rem**!

### License

GPLv3 - see [LICENSE](https://raw.githubusercontent.com/mborho/rem/master/LICENSE)

