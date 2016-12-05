package main

import (
	"golang.org/x/net/html"
	"strings"
	"net/http"
	"io/ioutil"
	"net/url"
)


// getUrlsFromPage parses the HTML content of the page passed in parameters
// and return an array of Pages.
//
// All Pages will contain an absolute URL.
//
// Href are filtered as follow :
//	- Should not start by '#' or 'mailto'
//	- Should not contain any 'action' parameters (used by Wikipedia forms)
//
func getUrlsFromPage(page *Page) []*Page {
	var pages []*Page

	resp, err := http.Get(page.Url.String())
	if err != nil {
		return pages
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pages
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return pages
	}

	// Retrieve and parse all hrefs from HTML body
	var f func(*html.Node)
	f = func(n*html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					newUrl, err := url.Parse(a.Val)
					if err != nil {
						continue
					}

					// Get rid of all the fragments
					if strings.HasPrefix(a.Val, "#") || strings.HasPrefix(a.Val, "mailto") {
						continue
					}
					// Bypass URL with action=edit
					if newUrl.Query().Get("action") != "" {
						continue
					}
					// Default to https for absolute URL starting with //
					if strings.HasPrefix(a.Val, "//") {
						newUrl.Scheme = "https"
					}
					// Build up absolute URL based on relative once
					if newUrl.Host == "" {
						newUrl.Host = page.Url.Host
						newUrl.Scheme = page.Url.Scheme
					}

					pages = append(pages, &Page{newUrl, page.Url.String()})
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return pages
}

