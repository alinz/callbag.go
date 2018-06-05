package main

import (
	"fmt"

	callbag "github.com/alinz/go-callbag"
)

func main() {
	src := callbag.Pipe(
		callbag.FromIter(1, 2, 3, 4),
	)

	callbag.Pipe(
		src,
		callbag.ForEach(func(val interface{}) {
			fmt.Println(val)
		}),
	)
}
