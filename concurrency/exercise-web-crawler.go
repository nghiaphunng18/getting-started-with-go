package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type SafeCache struct {
	mu      sync.Mutex
	visited map[string]bool
}

func (s *SafeCache) MarkVisited(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.visited[url] {
		return false
	}
	s.visited[url] = true
	return true
}

func Crawl(url string, depth int, fetcher Fetcher, cache *SafeCache, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	if !cache.MarkVisited(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, cache, wg)
	}
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

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
			"https://golang.org/pkg/",
		},
	},
}

func main() {
	var wg sync.WaitGroup
	cache := &SafeCache{visited: make(map[string]bool)}

	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, cache, &wg)
	wg.Wait()
	fmt.Println("\nDone")
}
