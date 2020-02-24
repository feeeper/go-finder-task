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

func main() {
	wg := &sync.WaitGroup{}
	ch := make(chan string, 5)
	scanner := bufio.NewScanner(os.Stdin)
	counter := 0
	mutex := new(sync.Mutex)
	regex := regexp.MustCompile("Go")

	go func() {
		for url := range ch {
			resp, err := http.Get(url)
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
			matches := regex.FindAllString(string(body), -1)
			count := len(matches)
			fmt.Println(url, "\t", count)

			mutex.Lock()
			counter += count
			mutex.Unlock()

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
