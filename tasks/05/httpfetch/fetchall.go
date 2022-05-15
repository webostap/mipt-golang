package httpfetch

import (
	"net/http"
	"sync"
)

type Request struct {
	Method string
	URL    string
	Body   []byte
}

type Result struct {
	StatusCode int
	Error      error
}

func FetchAll(c *http.Client, requests []Request) []Result {
	results := make([]Result, len(requests))
	var wg sync.WaitGroup
	for idx := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if resp, err := c.Get(requests[i].URL); err != nil {
				results[i].Error = err
			} else {
				results[i].StatusCode = resp.StatusCode
			}
		}(idx)
	}
	wg.Wait()
	return results
}
