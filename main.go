package main

import (
	"log"
	"sync"
	"time"
	"runtime"
	"flag"
	"net/url"
	"os"

	"github.com/PurpureGecko/go-lfc"
	"github.com/patrickmn/go-cache"
)


var (
	urlStart, urlDest string
	q *lfc.Queue
	nbWorkers, nbProcs int
	wg sync.WaitGroup
	pageCounter, maxVisitedPages int
	useWatcher bool
)


func main() {

	flag.StringVar(&urlStart, "start",  "", "Starting URL of the Wikirace")
	flag.StringVar(&urlDest, "dest", "", "Final destination of the Wikirace")
	flag.BoolVar(&useWatcher, "watcher", false, "If used, a watcher (goroutine) will be started to give live statistics")
	flag.IntVar(&nbWorkers, "workers", 10, "Number of workers (goroutines) to start")
	flag.IntVar(&nbProcs, "cpus", -1, "Override GOMAXPROCS if > 1")
	flag.IntVar(&maxVisitedPages, "max-visited-pages", -1, "If defined, will stop the process at the specified threshold (best effort). To be used with -watcher")

	flag.Parse()

	if nbProcs > 1 {
		runtime.GOMAXPROCS(nbProcs)
	}

	if urlStart == "" || urlDest == "" {
		log.Println("Use -start and -dest to set start and destination of the Wikirace")
		os.Exit(1)
	}

	// Create the Queue
	q = lfc.NewQueue()

	// Create the Cache
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	// Timer Start
	timerStart := time.Now()

	// Parse and enqueue start URL and end URL
	firstUrl, err := url.Parse(urlStart)
	if err != nil {
		panic(err)
	}
	q.Enqueue(&Page{firstUrl, ""})

	endUrl, err := url.Parse(urlDest)
	if err != nil {
		panic(err)
	}

	wg.Add(nbWorkers)

	done := make(chan struct{})
	for i := 0; i < nbWorkers; i++ {
		go worker(endUrl, q, c, &pageCounter, &wg, done)
	}

	if useWatcher {
		go watcher(q, &pageCounter, maxVisitedPages, done)
	}

	wg.Wait()

	// Timer End
	timerEnd := time.Since(timerStart)

	// Get path from cache
	path, found := retrievePathToPage(c, urlDest, urlStart)

	log.Printf("Results - Visited page(s) : %d - Size of the path : %d link(s) - Size of the queue : %d link(s) - Processing time : %s\n", pageCounter, len(path), q.Len(), timerEnd)
	if found {
		log.Printf("Path : %v\n", path)
	} else {
		log.Printf("Not path has been found !")
	}
}
