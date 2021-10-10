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
			regexp.MustCompile(`/media`),
			regexp.MustCompile(`dyncss`),
			regexp.MustCompile(`/favicon.ico`),
			regexp.MustCompile(`\.(gif|png|jpg|jpeg|svg)`),
		),
		colly.CacheDir("./.cache"),
		colly.Async(),
	)
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
	collector.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Println(link + " # css")
		e.Request.Visit(link)
	})
	collector.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		// fmt.Println(link + " # script")
		e.Request.Visit(link)
	})
	collector.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		// fmt.Println(link + " # img")
		e.Request.Visit(link)

	})
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Println(link)
		// Visit link found on page on a new thread
		e.Request.Visit(link)
	})
	collector.OnResponse(func(response *colly.Response) {
		if response.StatusCode == 200 {
			fmt.Println(response.Request.URL.String())
		}
	})
	collector.Visit("https://www.postgresql.org/docs/current/index.html")
	collector.Wait()
}
