package main

import (
	"fmt"
	"time"

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

	// callbag.Pipe(
	// 	callbag.FromIter(1, 2, 3, 4, 5),
	// 	callbag.Scan(func(prev interface{}, curr interface{}) interface{} {
	// 		a := prev.(int)
	// 		b := curr.(int)

	// 		return a + b
	// 	}, 0),
	// 	callbag.ForEach(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	//

	// callbag.Pipe(
	// 	callbag.Interval(1*time.Second),
	// 	callbag.Observe(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	//

	//

	// callbag.Pipe(
	// 	callbag.Interval(1*time.Second),
	// 	callbag.Take(5),
	// 	callbag.Observe(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	//

	// subject := callbag.Subject()

	// go func() {
	// 	callbag.Pipe(
	// 		subject,
	// 		callbag.Take(10),
	// 		callbag.Observe(func(val interface{}) {
	// 			fmt.Println("event: ", val)
	// 		}),
	// 	)
	// }()

	// time.Sleep(1 * time.Second)

	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		subject(callbag.NewData(fmt.Sprintf("Event %d", i)))
	// 	}
	// }()

	// go func() {
	// 	for i := 100; i < 200; i++ {
	// 		subject(callbag.NewData(fmt.Sprintf("Event %d", i)))
	// 	}
	// }()

	// time.Sleep(1 * time.Second)

	//

	// source := callbag.Concat(
	// 	callbag.FromIter(1, 2, 3, 4, 5),
	// 	callbag.FromIter("a", "b", "c"),
	// )

	// callbag.Pipe(
	// 	source,
	// 	callbag.ForEach(func(val interface{}) {
	// 		fmt.Println(val)
	// 	}),
	// )

	//

	source := callbag.Combine(
		callbag.Interval(100*time.Millisecond),
		callbag.Interval(350*time.Millisecond),
	)

	callbag.Pipe(
		source,
		callbag.Observe(func(val interface{}) {
			fmt.Println(val)
		}),
	)
}
