package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	ch := make(chan string, 5)
	scanner := bufio.NewScanner(os.Stdin)
	counter := 0
	mutex := new(sync.Mutex)

	go func() {
		for url := range ch {
			count, err := getOccurances(url)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(url, "\t", count)
				mutex.Lock()
				counter += count
				mutex.Unlock()
			}

			wg.Done()
		}
	}()

	for scanner.Scan() {
		wg.Add(1)

		url := scanner.Text()
		go func(url string) {
			ch <- url
		}(url)
	}

	wg.Wait()
	close(ch)

	fmt.Println("Total: ", "\t", counter)
}
