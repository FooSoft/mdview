# MdView

MdView is a simple tool for [GitHub Flavored Markdown](https://github.github.com/gfm/) files in your web browser. MdView
recursively renders the containing directory of the provided file path, and serves its contents over HTTP on port 8080.
The default browser is then launched, displaying the contents.

## Features

*   Works on Linux, Mac OS, and Windows.
*   Single executable; no installation or configuration required.
*   Preview is automatically refreshes when source files change.
*   Output supports GitHub styling, including light and dark modes.

## Usage

Simply execute MdView with the path of the Markdown file you wish to view:
```
$ mdview path/to/file.md
```
Execute MdView with `--help` to view options:
```
Usage mdview [options] [path]

Parameters:
  -port int
    	port (default 8080)
```

## Installation

If you already have the Go environment and toolchain set up, you can get the latest version by running:
```
$ go get github.com/FooSoft/mdview
```
Otherwise, you can use the [pre-built binaries](https://github.com/FooSoft/mdview/releases) from the project page.
