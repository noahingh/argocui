# ArgoCUI

[![Build Status](https://cloud.drone.io/api/badges/hanjunlee/argocui/status.svg)](https://cloud.drone.io/hanjunlee/argocui)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/hanjunlee/argocui/graphs/commit-activity)
[![GoDoc](https://godoc.org/github.com/hanjunlee/argocui?status.svg)](https://godoc.org/github.com/hanjunlee/argocui/pkg)


It support to manage Argo resource by CUI.

![ArgoCUI](./img/argocui.jpeg)

## Overview

The simple video clip for Argo CUI.

![ArgoCUI](./img/argocui-0.0.1.gif)

## Installation

### Release

You can download the binary file in [the release page](https://github.com/hanjunlee/argocui/releases).

### Brew

```shell
brew install argocui
```

### Source code

```shell
$ git clone git@github.com:hanjunlee/argocui.git
$ cd argocui
$ go build -o argocui ./cmd
```

## Command

```
Usage of acui  
  -debug
    	Debug mode.
  -trace
    	Debug as trace level.
  -ro
    	Read only mode. Some features such as terminate and delete doesn't work.
```

**Note that when you run the command it create the log file at** `$HOME/.argocui/log`.

## Keybinding

### List

 Key | Description
-----|-------------
 `/` | Set the search current view.
 `A` | Set as the global namespace, click again if you want to back.
 `k` | Move cursor up.
 `j` | Move cursor down.
 `H` | Move cursor up to the upper bound.
 `L` | Move cursor down to the bottom.
 `ctrl + n` | Switch to another namespace.
 `ctrl + l` | Display logs from Argo workflow.
 `ctrl + g` | Display the tree of Argo workflow.

### Search

 Key | Description
-----|-------------
 `enter` | Search Argo workflows which is matched with the pattern.
 `ctrl + u` | Clean.

### Logs & Tree

 Key | Description
-----|-------------
 `k` | Move cursor up.
 `j` | Move cursor down.
 `H` | Move cursor up to the upper bound.
 `L` | Move cursor down to the bottom.
 `esc` | Back to the list view.

## Changelog

[CHANGELOG.md](./docs/CHANGELOG.md)

## Contribute

[CONTRIBUTE.md](./docs/CONTRIBUTE.md)

