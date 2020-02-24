package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	ch := make(chan int)
	scanner := bufio.NewScanner(os.Stdin)

	i := 0
	for scanner.Scan() {
		wg.Add(1)
		url := scanner.Text()
		i += 1
		go func(address string) {
			resp, err := http.Get(url)
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
			fmt.Println(url, "\t", len(strings.Split(string(body), "Go"))-1)

			wg.Done()
		}(url)
	}

	wg.Wait()
	close(ch)
}
