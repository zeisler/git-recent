# git-recent

## About

Lists the most recent unmerged branches.

Uses auto complete to list branches so that the command is stored in history.

## Install

```shell
go install github.com/zeisler/git-recent@latest
git-recent -install -y
```

## Usage

**Use auto complete to list the last 20 unmerged branches.**

```shell
git-recent [TAB]
```

**With no arguments it will switch to the most recent branch.**

```shell
git-recent
```

## Future Improvements
* Make work as a git alias. Unable to get auto complete working with `git recent`, PR is welcome.