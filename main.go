package main

//@author Hansellll
//@date 1-10-26
//Golang website title aggregator

//Basic golang program that was made to showcase the power of concurrency and channels in golang
//Retrieves the first instance of a <title> on the home page of popular websites

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Article struct {
	URL   string
	Title string
	Body  string
}

func main() {
	var workers int
	fmt.Println("How many workers would you like to employ: ")
	fmt.Scanf("%d", &workers)
	start := time.Now()

	urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://reddit.com",
		"https://wikipedia.org",
		"https://medium.com",
		"https://dev.to",
		"https://hackerrank.com",
		"https://leetcode.com",
		"https://linkedin.com",
		"https://youtube.com",
		"https://netflix.com",
		"https://apple.com",
		"https://microsoft.com",
		"https://google.com",
		"https://facebook.com",
		"https://instagram.com",
		"https://twitch.tv",
		"https://discord.com",
		"https://slack.com",
		"https://zoom.us",
		"https://trello.com",
		"https://asana.com",
		"https://dropbox.com",
		"https://spotify.com",
		"https://soundcloud.com",
		"https://vimeo.com",
		"https://imgur.com",
		"https://flickr.com",
		"https://pinterest.com",
		"https://tumblr.com",
		"https://wordpress.com",
		"https://shopify.com",
		"https://etsy.com",
		"https://ebay.com",
		"https://craigslist.org",
		"https://yelp.com",
		"https://tripadvisor.com",
		"https://airbnb.com",
		"https://lyft.com",
		"https://doordash.com",
		"https://grubhub.com",
		"https://postmates.com",
	}

	// Build pipeline
	urlChan := generateUrls(urls)
	articleChan := fetchContent(urlChan, workers)
	titledChan := extractTitles(articleChan, workers)
	display(titledChan)

	fmt.Printf("Completed in %v\n", time.Since(start))
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//STAGE ONE

// Generate URLS
func generateUrls(urls []string) <-chan string {
	out := make(chan string)

	go func() {
		for _, url := range urls {
			out <- url
		}
		close(out)
	}() //Anonymous function call to the goroutine created above.
	return out
}

//STAGE TWO

// Fetch content from URLs

func fetchContent(urls <-chan string, workers int) <-chan Article {
	out := make(chan Article)
	var wg sync.WaitGroup

	// Launch worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() { // No parameter needed
			defer wg.Done()
			for url := range urls {
				// Just fetch, no worker ID needed
				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("Error fetching %s: %v\n", url, err)
					continue
				}
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				out <- Article{
					URL:  url,
					Body: string(body),
				}
			}
		}()
	}

	// Close output channel when all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

//STAGE THREE

// Find all titles in body of webpage using a worker pool
func extractTitles(articles <-chan Article, workers int) <-chan Article {
	out := make(chan Article)
	var wg sync.WaitGroup

	// Launch worker goroutines for title extraction
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for article := range articles {
				start := strings.Index(article.Body, "<title>")
				end := strings.Index(article.Body, "</title>")

				// Check that both tags exist AND are in correct order
				if start != -1 && end != -1 && start < end {
					article.Title = article.Body[start+7 : end]
				} else {
					article.Title = "No title found"
				}

				out <- article
			}
		}()
	}

	// Close output channel when all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Stage FOUR

// Display results
func display(articles <-chan Article) {
	for article := range articles {
		fmt.Printf("URL: %s\nTitle: %s\n\n", article.URL, article.Title)
	}
}
