package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
)

func getOccurances(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	regex := regexp.MustCompile("Go")
	matches := regex.FindAllString(string(body), -1)
	count := len(matches)

	return count, nil
}

func main() {
	wg := &sync.WaitGroup{}
	ch := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)
	counter := 0
	mutex := new(sync.Mutex)
	maxConcurrency := 2

	// ^C handling
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	go func() {
		for _ = range interrupts {
			fmt.Fprintf(os.Stdout, "Total: \t%d\n", counter)
			os.Exit(1)
		}
	}()

	// Exit on empty input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintf(os.Stderr, "Invalid input\n")
		os.Exit(1)
	}

	// Reading from stdin and process urls
	urlCount := 0
	for scanner.Scan() {
		url := scanner.Text()
		urlCount += 1
		if maxConcurrency >= urlCount {
			wg.Add(1)
			go func() {
				for url := range ch {
					count, err := getOccurances(url)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s\n", err.Error())
						fmt.Fprintf(os.Stdout, "%s\t0\n", url)
					} else {
						fmt.Fprintf(os.Stdout, "%s\t%d\n", url, count)
						mutex.Lock()
						counter += count
						mutex.Unlock()
					}
				}
				wg.Done()
			}()
		}

		ch <- url
	}

	close(ch)
	wg.Wait()
	fmt.Fprintf(os.Stdout, "Total: \t%d\n", counter)
}
