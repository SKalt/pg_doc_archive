package main

// note that relative paths should be interpreted relative to the repo root

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func main() {
	zip, err := gzip.NewWriterLevel(os.Stdout, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zip.Close()
	tapeArchive := tar.NewWriter(zip)
	defer tapeArchive.Close()

	pgLicense, err := os.Open("./licenses/postgres.license.md")
	if err != nil {
		log.Fatal(err)
	}
	info, err := pgLicense.Stat()
	if err != nil {
		log.Fatal(err)
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		log.Fatal(err)
	}
	header.Name = "/license.txt"
	err = tapeArchive.WriteHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(tapeArchive, pgLicense)
	if err != nil {
		log.Fatal(err)
	}
	pgLicense.Close()

	debugLog, err := os.CreateTemp(os.TempDir(), "archive.*.log")
	fmt.Fprintf(os.Stderr, "debug log at: %s\n", debugLog.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer debugLog.Close()

	collector := colly.NewCollector(
		colly.AllowedDomains("www.postgresql.org"),
		colly.Debugger(&debug.LogDebugger{Output: debugLog}),
		colly.URLFilters(
			regexp.MustCompile(`/docs/(current|\d.+)/`),
			regexp.MustCompile(`/about/license`),
			regexp.MustCompile(`/media`),
			regexp.MustCompile(`dyncss`),
			regexp.MustCompile(`/favicon.ico`),
			regexp.MustCompile(`\.(gif|png|jpg|jpeg|svg)`),
		),
		colly.CacheDir("./.cache"), // reuse cached site data
		// should be synchronous, one url at a time
	)

	visited := map[string]uint8{}
	collector.OnResponse(func(response *colly.Response) {
		if response.StatusCode != 200 {
			log.Fatalf("%+v", response)
		}
		path := "/" + strings.TrimLeft(response.Request.URL.Path, "/")
		// trim absolute url -> absolute path
		if strings.HasSuffix(path, "/") {
			path = path + "index.html"
		}
		if _, ok := visited[path]; !ok {
			visited[path] = 1
		} else {
			visited[path]++
			return
		}
		header := tar.Header{
			Name:   path,
			Size:   int64(len(response.Body)),
			Mode:   0666, // rw-rw-rw-
			Format: tar.FormatPAX,
		}
		err := tapeArchive.WriteHeader(&header)
		if err != nil {
			log.Fatal(err)
		}
		_, err = tapeArchive.Write(response.Body)
		if err != nil {
			log.Fatalf("%+v [%s; %d bytes]\n", err, path, len(response.Body))
		}
	})
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	bar := pb.Full.Start(len(lines))
	for _, url := range lines {
		_ = collector.Visit(url)
		bar.Increment()
	}
	bar.Finish()
}
