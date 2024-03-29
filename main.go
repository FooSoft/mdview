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

	"foosoft.net/projects/goldsmith"
	"foosoft.net/projects/goldsmith-components/devserver"
	"foosoft.net/projects/goldsmith-components/filters/operator"
	"foosoft.net/projects/goldsmith-components/filters/wildcard"
	"foosoft.net/projects/goldsmith-components/plugins/document"
	"foosoft.net/projects/goldsmith-components/plugins/frontmatter"
	"foosoft.net/projects/goldsmith-components/plugins/livejs"
	"foosoft.net/projects/goldsmith-components/plugins/markdown"
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

	allowedPaths := []string{
		"**/*.gif",
		"**/*.html",
		"**/*.jpeg",
		"**/*.jpg",
		"**/*.md",
		"**/*.png",
		"**/*.svg",
	}

	forbiddenPaths := []string{
		"**/.*/**",
	}

	errs := goldsmith.Begin(contentDir).
		Clean(true).
		FilterPush(wildcard.New(allowedPaths...)).
		FilterPush(operator.Not(wildcard.New(forbiddenPaths...))).
		Chain(frontmatter.New()).
		Chain(markdown.NewWithGoldmark(gm)).
		Chain(livejs.New()).
		Chain(document.New(embedCss)).
		End(buildDir)

	for _, err := range errs {
		log.Print(err)
	}

	if !self.open {
		url := fmt.Sprintf("http://127.0.0.1:%d/%s", self.port, self.path)
		log.Printf("opening %s in browser...", url)
		webbrowser.Open(url)
		self.open = true
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [options] [path]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Parameters:")
		flag.PrintDefaults()
	}

	port := flag.Int("port", 8080, "port")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := flag.Arg(0)
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	var contentName string
	contentDir := path
	if !info.IsDir() {
		contentName = filepath.Base(path)
		contentExt := filepath.Ext(contentName)
		switch contentExt {
		case ".md", ".markdown":
			contentName = strings.TrimSuffix(contentName, contentExt)
			contentName += ".html"
		}

		contentDir = filepath.Dir(path)
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
