package callbag

import (
	"reflect"
	"time"
)

func setInterval(fn func(), timeout time.Duration) func() {
	ticker := time.NewTicker(timeout)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fn()
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(clear)
	}
}

func Interval(period time.Duration) Source {
	return func(p Payload) {
		var sink Source
		var i int

		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			clear := setInterval(func() {
				sink(NewData(i))
				i++
			}, period)

			sink(NewGreets(func(p Payload) {
				if _, ok := p.(Terminate); ok {
					clear()
				}
			}))

		default:
			return
		}
	}
}

func PausableInterval(period time.Duration) Source {
	return func(p Payload) {
		var sink Source
		var i int
		var clear func()

		isClearNil := func() bool {
			return clear == nil || reflect.ValueOf(clear).IsNil()
		}

		resume := func() {
			clear = setInterval(func() {
				sink(NewData(i))
				i++
			}, period)
		}

		pause := func() {
			if !isClearNil() {
				clear()
				clear = nil
			}
		}

		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if isClearNil() {
						resume()
					} else {
						pause()
						clear = nil
					}
				case Terminate:
					if isClearNil() {
						pause()
					}
				}
			}))

			resume()

		default:
			return
		}
	}
}
