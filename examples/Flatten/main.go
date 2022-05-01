package main

import (
	"fmt"
	"sync"

	cb "github.com/alinz/callbag.go"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	cb.Pipe4(
		cb.FromRange(1, 11, 1),
		cb.Group[int](2),
		cb.Flatten[int](),
		cb.ForEach(func(value int, done bool) {
			if done {
				wg.Done()
				return
			}

			fmt.Printf("value is %v\n", value)
		}),
	)

	wg.Wait()
}
