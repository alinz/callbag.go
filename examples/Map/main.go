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
		cb.FromRange(1, 11, 1),
		cb.Map(func(value int) string {
			return fmt.Sprintf("Hello %d", value)
		}),
		cb.ForEach(func(value string, done bool) {
			if done {
				wg.Done()
				return
			}

			fmt.Printf("value is %s\n", value)
		}),
	)

	wg.Wait()
}
