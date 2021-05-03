package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	var t Task
	var wg sync.WaitGroup
	sem := make(chan bool, 10)
	for i := 1; i <= 100; i++ {
		// number of go routines are more than 10 in this case because there is main go routine, there
		// is garbage collector go routine
		fmt.Println(runtime.NumGoroutine())
		wg.Add(1)
		sem <- true
		go func(i int) {
			defer wg.Done()
			// we need to release the semaphore, or else we will execute only 10 go routines
			// all the other 90 go routines will be blocked forever
			defer func() { <-sem }()
			res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i))
			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Title %d is : %s \n", t.ID, t.Title)
		}(i)
	}
	wg.Wait()
}
