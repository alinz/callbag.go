# Callbag.go

this is an implementation for [Callbag](https://github.com/callbag/callbag)'s base core functionality in Go using new generics feature in 1.18.

At the moment the following functions have been implemented

- [x] Interval
- [x] FromSlice
- [x] FromRange
- [x] FromChannel
- [x] Filter
- [x] Map
- [x] Reduce
- [x] Take
- [x] ForEach
- [x] Pipe

# Installtion

```bash
go get github.com/alinz/callbag.go
```

# Example

count all numbers from 0 to 10

```go
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
```
