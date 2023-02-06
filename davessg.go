package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var verbose bool = false

type File struct {
	path    string
	ext     string
	modTime time.Time
	outPath string
}

func (f File) isSrcNewer() bool {
	destInfo, err := os.Stat(f.outPath)
	if err != nil {
		return true
	}
	return f.modTime.After(destInfo.ModTime())
}

func printf(f string, vars ...any) {
	fmt.Printf(f+"\n", vars...)
}

func debugf(f string, vars ...any) {
	if verbose {
		fmt.Printf(f+"\n", vars...)
	}
}

func findFiles(walkDir string, buildDir string) []File {
	var files = make([]File, 0, 100)
	filepath.WalkDir(walkDir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("Error walking content dir: %v", err)
		}

		if dir.IsDir() {
			return nil
		}

		file := File{}
		file.path = path

		fileInfo, err := dir.Info()
		if err != nil {
			log.Fatalf("Error reading file info for %v: %v", path, err)
		}

		file.modTime = fileInfo.ModTime()
		file.ext = strings.ToLower(filepath.Ext(path))
		file.outPath = mapDir(walkDir, path, file.ext, buildDir)
		files = append(files, file)

		return nil
	})

	return files
}

func mapDir(sourceDir string, srcPath string, ext string, buildDir string) string {
	var outPath = srcPath

	// Replace .md -> .html
	if ext == ".md" {
		if strings.ToLower(srcPath) == path.Join(sourceDir, "index.md") {
			// Special case for root index.md -> index.html (not a bundle)
			outPath = path.Join(buildDir, "index.html")
		} else {
			outPath = path.Join(outPath[:len(outPath)-len(ext)], "/index.html")
		}
	}

	// Content dir -> build dir prefix
	outPath = strings.Replace(outPath, sourceDir, buildDir, 1)
	debugf("mapping %v -> %v", srcPath, outPath)

	return outPath
}

func loadFile(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func build(files []File, baseurl string, force bool) {
	indexTmpl := loadFile(path.Join("templates", "index.html"))

	// Initialise goldmark
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	for _, file := range files {
		if file.isSrcNewer() || force {
			printf("Building %v", file.path)
			fileBytes := loadFile(file.path)

			if file.ext == ".md" {
				// Render markdown to html
				var buf bytes.Buffer
				err := md.Convert(fileBytes, &buf)
				if err != nil {
					log.Fatal(err)
				}
				fileBytes = buf.Bytes()
			}

			if file.ext == ".md" || file.ext == ".html" || file.ext == ".htm" || file.ext == ".js" {
				// Apply template
				fileBytes = bytes.ReplaceAll(indexTmpl, []byte("{{ $content }}"), fileBytes)

				// Do some string replacements
				fileBytes = bytes.ReplaceAll(fileBytes, []byte("{{ $baseurl }}"), []byte(baseurl))
			}

			// Write file
			os.MkdirAll(filepath.Dir(file.outPath), 0755)
			err := os.WriteFile(file.outPath, fileBytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			debugf("Skipping %v", file.path)
		}
	}
}

func main() {
	var serve bool
	var bindAddr string
	var force bool
	var buildDir string
	var sourceDir string
	var baseurl string

	flag.BoolVar(&serve, "serve", false, "Start web server")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.StringVar(&bindAddr, "bind-addr", "localhost:8009", "Listen address for web server (use with -serve)")
	flag.StringVar(&buildDir, "build-dir", "build/", "Build directory (created if necessary)")
	flag.StringVar(&sourceDir, "source-dir", "content/", "Source content dir")
	flag.StringVar(&baseurl, "base-url", "/", "Base URL")
	flag.BoolVar(&force, "force", false, "Overwrite existing build files, if necessary")
	flag.Parse()

	files := findFiles(sourceDir, buildDir)
	printf("Found %v content file(s)", len(files))

	staticFiles := findFiles(path.Join("templates", "static"), path.Join(buildDir, "static"))
	printf("Found %v static file(s)", len(staticFiles))

	files = append(files, staticFiles...)
	build(files, baseurl, force)

	if serve {
		printf("Listening on %v", bindAddr)
		log.Fatal(http.ListenAndServe(bindAddr, http.FileServer(http.Dir(buildDir))))
	}
}
