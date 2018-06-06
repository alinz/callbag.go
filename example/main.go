package main

import (
	"fmt"
	"time"

	callbag "github.com/alinz/go-callbag"
)

func main() {
	callbag.Pipe(
		callbag.Interval(1*time.Second),
		callbag.Map(func(val interface{}) interface{} {
			n := val.(int)
			return n + 1
		}),
		callbag.Filter(func(val interface{}) bool {
			n := val.(int)
			return n%2 == 0
		}),
		callbag.ForEach(func(val interface{}) {
			fmt.Println(val)
		}),
	)

	time.Sleep(12 * time.Second)
}
