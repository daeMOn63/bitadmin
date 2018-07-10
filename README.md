BitAdmin
===

BitAdmin is a command line tool aimed to help and speed up Atlassian Bitbucket (Self-Hosted version) administration.
It can be easily used as it is, or wrapped inside some bash scripts to avoid repetitive and long commands.

Note that it is not compatible with Bitbucket Cloud in current state.

<!-- toc -->
- [Overview](#overview)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Available Command List](#available-commands-list)
<!-- tocstop -->

## Overview

BitAdmin goal is to ease and speedup the administration of Atlassion BitBucket, when the users / projects / repositories database is getting big.
It aims to offer handy commands for common administration tasks (creating repository with predefined settings, settings user permissions...) by calling the BitBucket api (exposed from the [BitClient](https://github.com/daeMOn63/bitclient))

Also, it try to stay easily extensible, meaning adding a new custom command for a particular need must be quick and not require much boilerplate. Check for the [commands](commands) sources for more details.

## Installation
Make sure you have a working Go environment.  Go version 1.2+ is supported.  [See
the install instructions for Go](http://golang.org/doc/install.html).

To download bitadmin, simply run:
```
$ go get github.com/daeMOn63/bitadmin
```

Dependencies are managed using [Dep](https://github.com/golang/dep). Make sure to follow the [installation instructions](https://github.com/golang/dep#setup) to setup ```dep``` on your workstation.

To install the dependencies, run from the project root:
```
$ dep ensure
```

Then install the tool
```
$ go install
```

Make sure your `$GOPATH/bin` folder is in your path and the ```bitadmin``` command will be available from everywhere

### Autocompletion

Command line auto completion can save even more time while typing commands, as the tool provide autocomplete for :
- builtin commands
- commands arguments
- stash usernames
- stash project keys
- stash repository slugs
- permissions
- restrictions
- branchRefs

To enable it, add the following to your ~/.bashrc
```
PROG=bitadmin source $GOPATH/src/github.com/daeMOn63/bitadmin/vendor/github.com/urfave/cli/autocomplete/bash_autocomplete
```

And source it to make it effective
```
$ source ~/.bashrc
```

You will still need to warmup the cache for enabling stash fields completion, see [cache warmup](#cache-warmup) section for more details

### Cache warmup

Make sure to run the following command to preload the latest data from the BitBucket server, and make them available to autocompletion:

```
$ bitadmin --user YOUR_USERNAME --password ~/.bitadmin_secret --url http://stash.server.com cache warmup
```
You'll need to run this as often as you want to refresh your cached data with fresh one with the server.
But note that it's not needed by anything else than autocomplete.

## Getting started

The bitadmin binary provide built in documentation:
```
$ bitadmin
NAME:
   bitadmin - Bitbucket cli administration tool

USAGE:
   bitadmin [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Flavien Binet <https://github.com/daeMOn63/bitadmin>

COMMANDS:
     cache       Caching data for faster operation
     repository  Repository opertations
     user        User opertations
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --password <file>  Read password from <file>
   --url <url>        <url> of the bitbucket server
   --user <username>  Authenticate on bitbucket with <username>
   --help, -h         show help
   --version, -v      print the version
```

### Global options

User, password and url are all mandatory. First two should be filled with your BitBucket administrator credentials, and the url must point to your server.


To not expose your admin password over batch history or process list, it must be provided as a file.
You can create it like so :
```
$ set +o history
$ echo -n "password" > ~/.bitadmin_secret
$ chmod 600 ~/.bitadmin_secret
$ set -o history
```

File descriptor are also supported in case you don't want to write on disk, just remember to disable history first:
```
$ set +o history
$ bitadmin --password <(echo -n "mysecret") --user admin --url..
```

From here you should be good to go, try now
```
$ bitadmin --user YOUR_USERNAME --password ~/.bitadmin_secret --url http://stash.server.com cache warmup
```

And no errors should be reported.

Best might be to create an alias in ~/.bashrc to avoid repeating those settings all the time:
```
alias bitadmin='bitadmin --user YOUR_USERNAME --password ~/.bitadmin_secret --url "http://stash.server.com"'
```

## Available Commands List

This is the current supported commands tree:

```
- cache
    |- clear
    |- dump
    |- warmup
- repository
    |- create
    |- clone-settings
    |- set-branch-restriction
    |- set-branching-model
    |- set-pr-settings
    |- show-permission
    |- enable-sonar-cleanup
    |- set-default-reviewers
    |- move
- user
    |- grant
    |- unset-permissions
- group
    |- grant
- hooks
    |- list
    |- protect-unmerged-branch
        |- enable
        |- disable
    |- stash-eol-check
        |- enable
        |- disable
    |- yet-another-commit-checker
        |- enable
        |- disable
    |- reject-force-push
        |- enable
        |- disable
```

You can get more informations about a particular command or group by using the --help flag, available on everything :
```
$ bitadmin cache --help
$ bitadmin repository create --help
...
```
