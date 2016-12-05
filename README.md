# Wikiracer

Wikiracer is an implementation of the Wikipedia game : [Wikiracing](https://en.wikipedia.org/wiki/Wikipedia:Wikirace) using Go.

The main strategy is to parse the HTML content of the pages, stack all the found hrefs in a Queue (Lock-free queue), cache the visited links to retrieve the path and get to the final destination as fast as possible.


## Used libraries

I made the choice of using already existing libraries instead of reinventing the wheel. Here is what I'm using :

* Lock-Free Queue : [go-lfc](https://github.com/PurpureGecko/go-lfc)
* In-Memory Cache : [go-cache](https://github.com/patrickmn/go-cache)
  
## Usage

To build the Wikiracer :

```
make build
```

To play with it :

```
$ ./wikiracer -h
Usage of ./wikiracer:
  -cpus int
        Override GOMAXPROCS if > 1 (default -1)
  -dest string
        Final destination of the Wikirace
  -max-visited-pages int
        If defined, will stop the process at the specified threshold (best effort). To be used with -watcher (default -1)
  -start string
        Starting URL of the Wikirace
  -watcher
        If used, a watcher (goroutine) will be started to give live statistics
  -workers int
        Number of workers (goroutines) to start (default 10)
```

## Some results

```
$ ./wikiracer -cpus 3 -workers 2 -start 'https://en.wikipedia.org/wiki/Flash_(comics)' -dest 'http://www.dccomics.com/characters/reverse-flash'
2016/12/04 23:37:46 Results - Visited pages : 272 - Size of the path : 3 links - Size of the queue : 55436 links - Processing time : 16.882118982s
2016/12/04 23:37:46 Path : [https://en.wikipedia.org/wiki/Flash_(comics) http://www.dccomics.com/sites/theflash/ http://www.dccomics.com/characters/reverse-flash]
```

```
$ ./wikiracer -cpus 3 -workers 2 -start 'https://segment.com/blog/' -dest 'https://www.terraform.io/'
2016/12/04 23:31:53 Results - Visited pages : 40 - Size of the path : 3 links - Size of the queue : 387 links - Processing time : 1.621672197s
2016/12/04 23:31:53 Path : [https://segment.com/blog/ https://segment.com/blog/the-segment-aws-stack/ https://www.terraform.io/]
```

```
$ ./wikiracer -cpus 3 -workers 2 -start 'https://en.wikipedia.org/wiki/Walt_Disney' -dest 'https://en.wikipedia.org/wiki/Los_Angeles'
2016/12/04 23:21:57 Results - Visited pages : 5 - Size of the path : 3 links - Size of the queue : 6945 links - Processing time : 418.665527ms
2016/12/04 23:21:57 Path : [https://en.wikipedia.org/wiki/Walt_Disney https://en.wikipedia.org/wiki/The_Walt_Disney_Company https://en.wikipedia.org/wiki/Los_Angeles]
```

```
$ ./wikiracer -cpus 3 -workers 2 -start 'https://en.wikipedia.org/wiki/Linux' -dest 'https://en.wikipedia.org/wiki/Paris'
2016/12/04 23:25:29 Results - Visited pages : 436 - Size of the path : 3 links - Size of the queue : 74870 links - Processing time : 9.507102415s
2016/12/04 23:25:29 Path : [https://en.wikipedia.org/wiki/Linux https://en.wikipedia.org/wiki/GIMP https://en.wikipedia.org/wiki/Paris]
```

## To be continued

There is a lot of room for improvements :
 - The tool is currently lacking of unit testing
 - The current 'max-visited-pages' thresholds is a best effort limit since it is used within the 'watcher' gorouting. This goroutine is doing a 1sec sleep after each run.
 - Some benchmarking could be done to really test the speed of the solution
 - We could implement better interfaces for the Cache and the Queue in order to use external applications (Redis, Kafka, *MQ, etc)
