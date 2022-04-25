package main

import (
	"fmt"
	"sync"

	"github.com/alinz/callbag.go"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	src := callbag.FromRange(0, 11, 1) // range is [a, b), or upto 11

	sum := callbag.Reduce(func(base, value int) int {
		return base + value
	}, 0)

	sink := callbag.ForEach(func(value int, ok bool) {
		if !ok {
			wg.Done()
			return
		}

		fmt.Println("Result:", value)
	})

	callbag.Pipe3(src, sum, sink)

	wg.Wait()
}
