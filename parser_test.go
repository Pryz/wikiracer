package main

import(
	"testing"
	"net/url"
)

func BenchmarkGetUrlsFromPage(b *testing.B) {
	rawUrl := "https://en.wikipedia.org/wiki/Wikipedia:Wikirace"
	parsedUrl, _ := url.Parse(rawUrl)
	getUrlsFromPage(&Page{parsedUrl, ""})
}
