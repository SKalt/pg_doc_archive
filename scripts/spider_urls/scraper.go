package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func main() {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.postgresql.org"),
		colly.URLFilters(
			regexp.MustCompile(`/docs/(current|\d.+)/`),
			regexp.MustCompile(`/about/license`),
			regexp.MustCompile(`/media`),
			regexp.MustCompile(`dyncss`),
			regexp.MustCompile(`/favicon.ico`),
			regexp.MustCompile(`\.(gif|png|jpg|jpeg|svg)`),
		),
		colly.CacheDir("./.cache"),
		// colly.Async(),
	)
	// if err := collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4}); err != nil {
	// 	panic(err)
	// }
	collector.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		_ = collector.Visit(e.Request.AbsoluteURL(link))
	})
	collector.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		_ = collector.Visit(e.Request.AbsoluteURL(link))
	})
	collector.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		_ = collector.Visit(e.Request.AbsoluteURL(link))
	})
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		_ = collector.Visit(e.Request.AbsoluteURL(link))
	})
	collector.OnResponse(func(response *colly.Response) {
		if response.StatusCode == 200 {
			url := response.Request.URL.String()
			// if strings.HasSuffix(url, "/") {
			// url = url + "index.html"
			// }
			fmt.Println(url)
		}
	})
	if err := collector.Visit("https://www.postgresql.org/docs/current/index.html"); err != nil {
		panic(err)
	}
	collector.Wait()
}
