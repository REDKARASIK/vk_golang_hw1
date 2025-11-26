package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	keepAlive   = flag.Bool("k", false, "Enable HTTP KeepAlive")
	concurrency = flag.Int("c", 1, "Number of concurrent requests")
	requests    = flag.Int("n", 1, "Number of total requests to perform")
)

type result struct {
	duration time.Duration
	status   int
	error    bool
}

func main() {
	flag.Parse()
	url := flag.Arg(0)

	if url == "" {
		fmt.Println("Usage: ab -k -c <concurrency> -n <requests> <url>")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: !*keepAlive,
		},
		Timeout: 30 * time.Second,
	}

	results := make(chan result, *requests)
	var wg sync.WaitGroup
	sem := make(chan struct{}, *concurrency)

	// Atomic counter for progress reporting
	var completed int64

	start := time.Now()

	// Progress reporter
	progressDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				current := atomic.LoadInt64(&completed)
				fmt.Printf("Completed %d requests\n", current)
			case <-progressDone:
				// Print final progress if we didn't just print it
				current := atomic.LoadInt64(&completed)
				if current > 0 && current%10000 != 0 {
					fmt.Printf("Completed %d requests\n", current)
				}
				return
			}
		}
	}()

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			reqStart := time.Now()
			resp, err := client.Get(url)
			if err != nil {
				atomic.AddInt64(&completed, 1)
				results <- result{error: true}
				return
			}
			defer resp.Body.Close()

			io.Copy(io.Discard, resp.Body)
			atomic.AddInt64(&completed, 1)
			results <- result{
				duration: time.Since(reqStart),
				status:   resp.StatusCode,
			}
		}()
	}

	wg.Wait()
	close(progressDone) // Stop progress reporter
	close(results)
	totalTime := time.Since(start)

	var totalDuration time.Duration
	successful := 0
	failed := 0
	statusCodes := make(map[int]int)

	for r := range results {
		if r.error {
			failed++
			continue
		}
		successful++
		totalDuration += r.duration
		statusCodes[r.status]++
	}

	fmt.Printf("\nBenchmark results:\n")
	fmt.Printf("Total time: %.3fs\n", totalTime.Seconds())
	fmt.Printf("Requests: %d total, %d successful, %d failed\n",
		*requests, successful, failed)
	if successful > 0 {
		fmt.Printf("Average response time: %.3fs\n",
			totalDuration.Seconds()/float64(successful))
	}
	fmt.Printf("Requests per second: %.2f\n",
		float64(*requests)/totalTime.Seconds())
	fmt.Println("Status code distribution:")
	for code, count := range statusCodes {
		fmt.Printf("  %d: %d responses\n", code, count)
	}
}
