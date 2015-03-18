# Bob
Build services to Jenkins

## Installation guide for users who are not familiar with golang

step 1: install go on your machine and setup the environment by following the instructions from the following link:

https://golang.org/doc/install

step 2: make a directory in `$HOME` where all your Go code will live
```
mkdir ~/gocode
mkdir ~/gocode/src
mkdir ~/gocode/bin
mkdir ~/gocode/pkg
```

step 3: Tell Go to use that as your GOPATH: `export GOPATH=~/gocode`

step 4: Save your GOPATH so that it will apply to all future shell sessions: 

`echo export GOPATH=$GOPATH >> ~/.bash_profile`

step 5: You may need to install Mercurial by following the instructions below

http://mercurial.selenic.com/downloads

step 6: Get the repo 
either with `go get github.com/ziyadparekh/Bob` or
```
cd $GOPATH/src/github.com/ziyadparekh
git clone git@github.com:ziyadparekh/Bob.git
go get -u
go install
```

step 7: Ensure the `$GOPATH/bin` directory is in your PATH so that you can reference bob from anywhere.
```
export PATH="$PATH:$GOPATH/bin"
echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bash_profile
```

Lastly, let’s verify that it works: :)

```
$ bob help
~
NAME:
   bob - Build Services to Deathstar

USAGE:
   bob [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
  Ziyad Parekh - <unknown@email>

COMMANDS:
   build	Build a service to deathstar!
   list		List all the services that are buildable
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
   ```
