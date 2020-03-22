![ArgoCUI](./img/argocui.jpeg)

# ArgoCUI - CLI To Manage Argo resource.

[![Build Status](https://cloud.drone.io/api/badges/hanjunlee/argocui/status.svg)](https://cloud.drone.io/hanjunlee/argocui) [![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/hanjunlee/argocui/graphs/commit-activity) [![GoDoc](https://godoc.org/github.com/hanjunlee/argocui?status.svg)](https://godoc.org/github.com/hanjunlee/argocui/pkg)

---

Argocui provides a terminal UI to manage [Argo](https://github.com/argoproj/argo) resources. The aim of this project is to deal with Argo resources such as `Workflow` and `CronWorkflow` as the `argo` command provides. 


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
Usage of argocui  
  -debug
    	Debug mode.
  -trace
    	Debug as trace level.
  -version
    	Check the version.
```

**Note that when you run the command it create the log file at** `$HOME/.argocui/log`.

## Keybinding

Command | Description 
--------|-------------
`H`     | Move the cursor to the top of a view.
`k`     | Move the cursor up.
`j`     | Move the cursor down.
`L`     | Move the cursor to the bottom of a view.
`esc`   | Back to the main view which displays the list of resources.
`:`     | Switch the [Kind](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#types-kinds) to another. It has `ns` i.e `Namespace` and `wf` i.e `Workflow` at this moment. For example, `:ns` switch the Kind into the `Namespace`. 
`/`     | Search resources which is matched with the pattern. 
`ctrl+g`| Display the detail of resource. It works for `Workflow`.
`ctrl+l`| Follow logs of resource. It works for `Workflow`.
`ctrl+delete` | Delete the resource. It works for `Workflow`.

## Changelog

[CHANGELOG.md](./docs/CHANGELOG.md)

## Contribute

[CONTRIBUTE.md](./docs/CONTRIBUTE.md)

## Special thanks to 

I owe a huge thanks to maintainer and contributers of [K9S](https://github.com/derailed/k9s).
