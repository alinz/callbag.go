package main

import (
	"fmt"
	"sync"

	cb "github.com/alinz/callbag.go"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	cb.Pipe3(
		cb.FromRange(1, 100, 1),
		cb.Filter(func(value int) bool {
			return value%2 == 0
		}),
		cb.ForEach(func(value int, done bool) {
			if done {
				wg.Done()
				return
			}

			// prints out all even number between [1,100)
			fmt.Printf("value is %d\n", value)
		}),
	)

	wg.Wait()
}
