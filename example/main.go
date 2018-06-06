package main

import (
	"fmt"

	callbag "github.com/alinz/go-callbag"
)

func main() {
	callbag.Pipe(
		callbag.FromIter(1, 2, 3, 4),
		callbag.Map(func(val interface{}) interface{} {
			n := val.(int)
			return n + 1
		}),
		callbag.ForEach(func(val interface{}) {
			fmt.Println(val)
		}),
	)
}
