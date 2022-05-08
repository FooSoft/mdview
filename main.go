package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-components/devserver"
	"github.com/FooSoft/goldsmith-components/plugins/document"
	"github.com/FooSoft/goldsmith-components/plugins/livejs"
	"github.com/FooSoft/goldsmith-components/plugins/markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/toqueteos/webbrowser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed css/github-markdown.css
var githubStyle string

//go:embed css/github-fixup.css
var githubFixup string

type builder struct {
	port int
	path string
	open bool
}

func embedCss(file *goldsmith.File, doc *goquery.Document) error {
	var styleBuilder strings.Builder
	styleBuilder.WriteString("<style type=\"text/css\">\n")
	styleBuilder.WriteString(githubStyle)
	styleBuilder.WriteString(githubFixup)
	styleBuilder.WriteString("</style>")

	doc.Find("body").AddClass("markdown-body")
	doc.Find("head").SetHtml(styleBuilder.String())

	return nil
}

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
		Chain(document.New(embedCss)).
		End(buildDir)

	for _, err := range errs {
		log.Print(err)
	}

	if !self.open {
		webbrowser.Open(fmt.Sprintf("http://127.0.0.1:%d/%s", self.port, self.path))
		self.open = true
	}
}

func main() {
	port := flag.Int("port", 8080, "port")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("unexpected number of arguments")
	}

	requestPath := flag.Arg(0)
	info, err := os.Stat(requestPath)
	if err != nil {
		log.Fatal(err)
	}

	var contentName string
	contentDir := requestPath
	if !info.IsDir() {
		contentName = filepath.Base(requestPath)
		contentExt := filepath.Ext(contentName)
		switch contentExt {
		case ".md", ".markdown":
			contentName = strings.TrimSuffix(contentName, contentExt)
			contentName += ".html"
		}

		contentDir = filepath.Dir(requestPath)
	}

	buildDir, err := ioutil.TempDir("", "mvd-*")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Println("cleaning up...")
		if err := os.RemoveAll(buildDir); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		b := &builder{port: *port, path: contentName}
		devserver.DevServe(b, *port, contentDir, buildDir, "")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
