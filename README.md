# rem - remember

rem is a small tool for remembering things on the command line. It can be used to store commands and execute them later on. Or to simply store lines of text.

### Build & Run

```sh
$ git clone https://github.com/mborho/rem
$ cd rem
$ GOPATH=`pwd` go build -o rem main.go
$ ./rem
```

### Statically linked binary 

[LINUX/x86](https://raw.githubusercontent.com/mborho/rem/master/bin/linux_amd64/rem) (2MB)

[DARWIN/386](https://raw.githubusercontent.com/mborho/rem/master/bin/linux_darwin/rem) (1.6MB)


### PPA ##
rem is available for Ubuntu x86_64 [ppa:martin-borho/rem](https://launchpad.net/~martin-borho/+archive/ubuntu/rem)


```
#!bash

sudo add-apt-repository ppa:martin-borho/rem
sudo apt-get update
sudo apt-get install rem
```

### Usage
```sh
$ rem command [argument]
```
### Commands

*    help - Shows this help.
*    add [string] - Adds a command/text.
*    rm [index] - Removes line with given index number.
*    echo [index] - Displays line with given index number.
*    here - Creates a .rem file in the given directory. Default: **~/.rem**
*    clear - Clears currently active .rem file, **./.rem** or **~/.rem**
*    [index] - Executes line with given index number.

Run **rem** without any arguments to list all stored commands/strings.

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

