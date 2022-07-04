<!-- +++
Area = "projects"
GitHub = "mdview"
Layout = "page"
Tags = ["generator", "golang", "markdown", "mdview", "mit license", "web"]
Description = "Tool to view Github Flavored Markdown files in your web browser."
Collection = "ProjectsActive"
+++ -->

# MdView

MdView is a tool for viewing [GitHub Flavored Markdown](https://github.github.com/gfm/) files in your web browser.
MdView recursively renders the parent directory of the user-provided file and serves its contents over HTTP on port
8080. The default browser is then launched, displaying the Markdown file as an HTML document.

## Features

*   Works on Linux, Mac OS, and Windows.
*   Single executable; no installation or configuration required.
*   Preview is automatically refreshed when source files change.
*   Output supports GitHub styles, including light and dark modes.
*   Automatically skips any included front matter header data.

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
$ go install foosoft.net/projects/mdview@latest
```
Otherwise, you can use the [pre-built binaries](https://github.com/FooSoft/mdview/releases) from the project page.
