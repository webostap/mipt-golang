package httpfetch2

import (
	"bytes"
	"context"
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

func FetchAll(ctx context.Context, c *http.Client, requests <-chan Request) <-chan Result {
	ch := make(chan Result)
	var wg sync.WaitGroup
	go func() {
		defer func() {
			wg.Wait()
			close(ch)
		}()
		for {
			select {
			case request, ok := <-requests:
				if ok {
					wg.Add(1)
					go func(req Request) {
						defer wg.Done()
						rr, _ := http.NewRequestWithContext(ctx, req.Method, req.URL, bytes.NewReader(req.Body))
						var res Result
						if resp, err := c.Do(rr); resp != nil {
							defer resp.Body.Close()
							res.StatusCode = resp.StatusCode
						} else {
							res.Error = err
						}
						ch <- res
					}(request)
				} else {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}
