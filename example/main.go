package main

import (
	"fmt"

	"github.com/alinz/go-callbag"
)

func main() {
	// callbag.Pipe(
	// 	callbag.Interval(1*time.Second),
	// 	callbag.Map(func(val interface{}) interface{} {
	// 		n := val.(int)
	// 		return n + 1
	// 	}),
	// 	callbag.Filter(func(val interface{}) bool {
	// 		n := val.(int)
	// 		return n%2 == 0
	// 	}),
	// 	callbag.ForEach(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	//

	// source := callbag.Merge(
	// 	callbag.Interval(100*time.Millisecond),
	// 	callbag.Interval(200*time.Millisecond),
	// )

	// callbag.Pipe(
	// 	source,
	// 	callbag.ForEach(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	// time.Sleep(12 * time.Second)

	//

	// callbag.Pipe(
	// 	callbag.FromIter(1, 2, 3, 4),
	// 	callbag.Map(func(val1 interface{}) interface{} {
	// 		return callbag.Pipe(
	// 			callbag.FromIter(5, 6, 7, 8),
	// 			callbag.Map(func(val2 interface{}) interface{} {
	// 				return fmt.Sprintf("%d%d", val1, val2)
	// 			}),
	// 		)
	// 	}),
	// 	callbag.Flatten(),
	// 	callbag.ForEach(func(val interface{}) {
	// 		fmt.Printf("%v\n", val)
	// 	}),
	// )

	//

	callbag.Pipe(
		callbag.FromIter(1, 2, 3, 4, 5),
		callbag.Scan(func(prev interface{}, curr interface{}) interface{} {
			a := prev.(int)
			b := curr.(int)

			return a + b
		}, 0),
		callbag.ForEach(func(val interface{}) {
			fmt.Println(val)
		}),
	)

}
