package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var fetched map[string]struct{}
var fetchedMu sync.Mutex

var wg sync.WaitGroup

var gorutineCount atomic.Int32

func init() {
	fetched = make(map[string]struct{})
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// Capture a counter for parallel debug
	id := gorutineCount.Add(1)
	fmt.Printf("[%d] starting new goroutine ...\n", id)

	// Once done, mark work as completed.
	defer wg.Done()

	// Max depth reached, return
	if depth <= 0 {
		fmt.Printf("[%d] max-depth reached: %v\n", id, depth)
		return
	}

	// Don't fetch if already fetched
	fetchedMu.Lock()
	if _, seen := fetched[url]; seen {
		fetchedMu.Unlock()
		fmt.Printf("[%d] already seen: %v\n", id, url)
		return
	}

	// Fetches the page and store this one as completed.
	body, urls, err := fetcher.Fetch(url)
	fetched[url] = struct{}{}
	fetchedMu.Unlock()

	// If error, report and return.
	if err != nil {
		fmt.Printf("[%d] error: %v\n", id, err)
		return
	}

	// Otherwise, let's move on and spawn more routines in parallel.
	fmt.Printf("[%d] fetched: %s %q\n", id, url, body)
	wg.Add(len(urls))
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher)
	}
}

func main() {
	// Spawn the first parallel task
	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher)

	// Wait until all work is done.
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	// Fake latency
	time.Sleep(123 * time.Millisecond)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
