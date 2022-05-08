package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-components/devserver"
	"github.com/FooSoft/goldsmith-components/plugins/livejs"
	"github.com/FooSoft/goldsmith-components/plugins/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type builder struct{}

func (self *builder) Build(contentDir, buildDir, cacheDir string) {
	log.Print("building...")

	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	errs := goldsmith.Begin(contentDir).
		Clean(true).
		Chain(markdown.NewWithGoldmark(gm)).
		Chain(livejs.New()).
		End(buildDir)

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	port := flag.Int("port", 8080, "port")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("unexpected number of arguments")
	}

	requestPath := flag.Arg(0)
	buildDir, err := ioutil.TempDir("", "mvd-*")
	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(requestPath)
	if err != nil {
		log.Fatal(err)
	}

	contentDir := requestPath
	if !info.IsDir() {
		contentDir = filepath.Dir(requestPath)
	}

	b := new(builder)
	devserver.DevServe(b, *port, contentDir, buildDir, "")
}
