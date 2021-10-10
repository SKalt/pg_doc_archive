package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

var logger = log.Default()

func main() {
	zip, err := gzip.NewWriterLevel(os.Stdout, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zip.Close()
	tapeArchive := tar.NewWriter(zip)
	defer tapeArchive.Close()

	collector := colly.NewCollector(
		colly.AllowedDomains("www.postgresql.org"),
		colly.URLFilters(
			regexp.MustCompile(`/docs/(current|\d.+)/`),
			regexp.MustCompile(`/media`),
			regexp.MustCompile(`dyncss`),
			regexp.MustCompile(`/favicon.ico`),
			regexp.MustCompile(`\.(gif|png|jpg|jpeg|svg)`),
		),
		colly.CacheDir("./.cache"), // reuse cached site data
		// should be sychronous, one url at a time
	)
	collector.OnResponse(func(response *colly.Response) {
		path := response.Request.URL.Path
		if path[0:2] == "//" { // trim absolute url -> absolute path
			path = path[1:]
		}
		logger.Printf(
			"visiting: %s [%d]\n",
			path, response.StatusCode)
		if response.StatusCode != 200 {
			return
		}

		header := tar.Header{
			Name:   response.Request.URL.Path,
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
			log.Fatal(err)
		}
	})
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		path := scanner.Text()
		logger.Printf("input: %s", path)
		collector.Visit("https://www.postgresql.org/" + path)
	}
}
