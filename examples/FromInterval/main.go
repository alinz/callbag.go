package main

import (
	"fmt"
	"sync"
	"time"

	cb "github.com/alinz/callbag.go"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	cb.Pipe2(
		cb.FromInterval[int](1*time.Second),
		cb.ForEach(func(value int, done bool) {
			if done {
				wg.Done()
				return
			}

			fmt.Printf("value is %d\n", value)
		}),
	)

	wg.Wait()
}
