package main

import (
	"log"
	"sync"
	"time"
	"net/url"

	"github.com/PurpureGecko/go-lfc"
	"github.com/patrickmn/go-cache"
)

// Page is the data structure used within the Queue
type Page struct {
	Url *url.URL
	ParentUrl string
}


func worker(endUrl *url.URL, q *lfc.Queue, c *cache.Cache, pCounter *int, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	for {
		select {
			case <- done:
				// A worker has found endUrl
				return
			default:
				// Keep going
		}
		nextPage, ok := q.Dequeue()
		if ok {
			*pCounter++
			curPage := nextPage.(*Page)

			if curPage.Url.String() == endUrl.String() {
				close(done)
				return
			}

			pages := getUrlsFromPage(curPage)
			for _, p := range(pages) {
				// If not already in the cache, process
				if _, found := c.Get(p.Url.String()); found {
					continue
				}
				c.Set(p.Url.String(), curPage.Url.String(), cache.NoExpiration)
				if p.Url.String() == endUrl.String() {
					close(done)
					return
				}
				q.Enqueue(p)
			}
		}
	}
}


func watcher(q *lfc.Queue, c *int, maxPages int ,done chan struct{}) {
	for {
		select {
			case <- done:
				log.Println("WATCHER - A worker found the link ! huhu !")
				return
			default:
				log.Printf("WATCHER - Queue size : %d\n", q.Len())
				log.Printf("WATCHER - Visited pages : %d\n", *c)
				if q.Len() > 0 {
					head := q.Get(1)
					last := head[0].(*Page)
					log.Printf("WATCHER - Queue head : %s\n", last.Url.String())
				}
				if maxPages > 0 && *c > maxPages {
					log.Printf("WATCHER - Maximum Visited pages reached : %d. No path has been found. Stopping the process...\n", maxPages)
					close(done)
				}
				time.Sleep(1*time.Second)
		}
	}
}


func retrievePathToPage(c *cache.Cache, page string, root string) ([]string, bool) {
	var path []string
	node, found := c.Get(page)
	if found == false {
		return nil, false
	}
	path = append(path, page)
	path = append(path, node.(string))
	for node != root {
		node, found = c.Get(node.(string))
		path = append(path, node.(string))
	}
	return reversePath(path), true
}


func reversePath(path []string) []string {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
