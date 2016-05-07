package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	// Can be a locally running server in case of no internet connection
	// In my case it's just a dummy server responding with 200 to get requests
	// sourceURLs = []string{
	// 	"http://localhost:8080",
	// }

	// Or a list of actual urls
	sourceURLs = []string{
		"http://www.google.com",
		"http://www.yahoo.com",
		"http://www.amazon.com",
	}
)

func main() {
	// Create an entropy generator based on duration of web requests
	wreg := NewWebRequestsEntropyGenerator(20, 200*time.Millisecond, sourceURLs)
	outputChannel := wreg.StartGenerating()

	// Start feeding random byte values off of outputChannel
	// and output them to stdout
	go func() {
		for randomByte := range *outputChannel {
			fmt.Printf("%c", randomByte)
		}
	}()

	fmt.Scanln()
}

type WebRequestsEntropyGenerator struct {
	Concurrency      int
	Sources          []string
	ThrottlingPeriod time.Duration
}

func NewWebRequestsEntropyGenerator(concurrency int, throttlingPeriod time.Duration, sources []string) *WebRequestsEntropyGenerator {
	// Limit concurrency
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > 20 {
		concurrency = 20
	}

	wreg := WebRequestsEntropyGenerator{
		Concurrency:      concurrency,
		Sources:          sources,
		ThrottlingPeriod: throttlingPeriod,
	}

	return &wreg
}

func (wreg *WebRequestsEntropyGenerator) StartGenerating() *chan byte {
	randomBytesChannel := make(chan byte, wreg.Concurrency)

	// start workers
	for i := 0; i < wreg.Concurrency; i++ {
		go func() {
			currentIdx := 0

			for {
				d, err := getRequestTime(wreg.Sources[currentIdx])

				// round-robin through sourceURLs
				currentIdx = (currentIdx + 1) % len(wreg.Sources)

				if err != nil {
					continue
				}

				randomBytesChannel <- byte(d % 255)

				time.Sleep(wreg.ThrottlingPeriod)
			}
		}()
	}

	return &randomBytesChannel
}

// Time a GET request to ns resolution
func getRequestTime(sourceURL string) (int64, error) {
	start := time.Now()

	_, err := http.Get(sourceURL)
	if err != nil {
		return 0, err
	}

	return time.Since(start).Nanoseconds(), nil
}
